package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	_ "modernc.org/sqlite"
)

// Структуры данных
type NewsItem struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Summary     string    `json:"summary"`
	Source      string    `json:"source"`
	URL         string    `json:"url"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	Category    string    `json:"category"`
}

type User struct {
	ID        int       `json:"id"`
	UUID      string    `json:"uuid"`
	Sources   []string  `json:"sources"`
	CreatedAt time.Time `json:"created_at"`
}

type SummarizeRequest struct {
	Content string `json:"content"`
}

type SummarizeResponse struct {
	Summary string `json:"summary"`
}

// Глобальная переменная для БД
var db *sql.DB

// Middleware для логирования запросов
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Запрос: %s %s", r.Method, r.URL.Path)
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Паника в обработчике: %v", err)
				http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Инициализация БД
	initDB()

	// Создание роутера
	r := mux.NewRouter()

	// Добавляем middleware для логирования
	r.Use(loggingMiddleware)

	// API эндпоинты
	api := r.PathPrefix("/api").Subrouter()

	// Пользователи
	api.HandleFunc("/users", createUser).Methods("POST")
	api.HandleFunc("/users/{uuid}/sources", updateUserSources).Methods("PUT")
	api.HandleFunc("/users/{uuid}/sources", getUserSources).Methods("GET")

	// Новости
	api.HandleFunc("/news", getNews).Methods("GET")
	api.HandleFunc("/news/{id}", getNewsItem).Methods("GET")
	api.HandleFunc("/news", addNews).Methods("POST")

	// ML эндпоинты
	api.HandleFunc("/summarize", summarizeText).Methods("POST")

	// CORS настройки
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	fmt.Printf("Сервер запущен на порту %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
} // Инициализация базы данных
func initDB() {
	var err error
	log.Println("Инициализация базы данных...")
	db, err = sql.Open("sqlite", "./news.db")
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	// Проверяем соединение
	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка проверки соединения с БД:", err)
	}
	log.Println("Соединение с БД установлено")

	// Создание таблиц
	createTables()
	log.Println("База данных готова к работе")
}

func createTables() {
	// Таблица пользователей
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT UNIQUE NOT NULL,
		sources TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Таблица новостей
	newsTable := `
	CREATE TABLE IF NOT EXISTS news (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		summary TEXT,
		source TEXT NOT NULL,
		url TEXT UNIQUE NOT NULL,
		published_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		category TEXT
	);`

	_, err := db.Exec(userTable)
	if err != nil {
		log.Fatal("Ошибка создания таблицы users:", err)
	}

	_, err = db.Exec(newsTable)
	if err != nil {
		log.Fatal("Ошибка создания таблицы news:", err)
	}
}

// Обработчики API

// Создание пользователя
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	sourcesJSON, _ := json.Marshal(user.Sources)

	query := "INSERT INTO users (uuid, sources) VALUES (?, ?)"
	result, err := db.Exec(query, user.UUID, string(sourcesJSON))
	if err != nil {
		http.Error(w, "Ошибка создания пользователя", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	user.ID = int(id)
	user.CreatedAt = time.Now()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Обновление источников пользователя
func updateUserSources(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	var sources []string
	err := json.NewDecoder(r.Body).Decode(&sources)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	sourcesJSON, _ := json.Marshal(sources)

	query := "UPDATE users SET sources = ? WHERE uuid = ?"
	_, err = db.Exec(query, string(sourcesJSON), uuid)
	if err != nil {
		http.Error(w, "Ошибка обновления источников", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// Получение источников пользователя
func getUserSources(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	var sourcesJSON string
	query := "SELECT sources FROM users WHERE uuid = ?"
	err := db.QueryRow(query, uuid).Scan(&sourcesJSON)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	var sources []string
	json.Unmarshal([]byte(sourcesJSON), &sources)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sources)
}

// Получение списка новостей
func getNews(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос на новости: %s", r.URL.String())
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "20"
	}

	source := r.URL.Query().Get("source")
	category := r.URL.Query().Get("category")

	query := "SELECT id, title, content, summary, source, url, published_at, created_at, category FROM news WHERE 1=1"
	args := []interface{}{}

	if source != "" {
		query += " AND source = ?"
		args = append(args, source)
	}

	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}

	query += " ORDER BY published_at DESC LIMIT ?"
	limitInt, _ := strconv.Atoi(limit)
	args = append(args, limitInt)

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, "Ошибка получения новостей", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var news []NewsItem
	for rows.Next() {
		var item NewsItem
		err := rows.Scan(&item.ID, &item.Title, &item.Content, &item.Summary,
			&item.Source, &item.URL, &item.PublishedAt, &item.CreatedAt, &item.Category)
		if err != nil {
			continue
		}
		news = append(news, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}

// Получение конкретной новости
func getNewsItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var item NewsItem
	query := "SELECT id, title, content, summary, source, url, published_at, created_at, category FROM news WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&item.ID, &item.Title, &item.Content,
		&item.Summary, &item.Source, &item.URL, &item.PublishedAt, &item.CreatedAt, &item.Category)
	if err != nil {
		http.Error(w, "Новость не найдена", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// Добавление новости
func addNews(w http.ResponseWriter, r *http.Request) {
	var item NewsItem
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO news (title, content, summary, source, url, published_at, category) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, item.Title, item.Content, item.Summary,
		item.Source, item.URL, item.PublishedAt, item.Category)
	if err != nil {
		http.Error(w, "Ошибка добавления новости", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	item.ID = int(id)
	item.CreatedAt = time.Now()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// Суммаризация текста через ML модуль
func summarizeText(w http.ResponseWriter, r *http.Request) {
	var req SummarizeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Вызов Python скрипта для суммаризации
	summary, err := callPythonSummarizer(req.Content)
	if err != nil {
		log.Printf("Ошибка суммаризации: %v", err)
		// Fallback - простое сокращение текста
		summary = req.Content[:min(200, len(req.Content))] + "..."
	}

	response := SummarizeResponse{Summary: summary}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Вызов Python скрипта для суммаризации
func callPythonSummarizer(content string) (string, error) {
	// Подготавливаем входные данные для Python скрипта
	input := map[string]interface{}{
		"content":        content,
		"sentence_count": 2,
	}

	inputJSON, err := json.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("ошибка подготовки данных: %v", err)
	}

	// Выполняем Python скрипт
	cmd := exec.Command("python", "../ml/summarize.py")
	cmd.Stdin = strings.NewReader(string(inputJSON))

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения Python скрипта: %v", err)
	}

	// Парсим результат
	var result map[string]interface{}
	err = json.Unmarshal(output, &result)
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга результата: %v", err)
	}

	if status, ok := result["status"].(string); ok && status == "error" {
		return "", fmt.Errorf("ошибка в Python скрипте: %v", result["error"])
	}

	summary, ok := result["summary"].(string)
	if !ok {
		return "", fmt.Errorf("неверный формат ответа от Python скрипта")
	}

	return summary, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
