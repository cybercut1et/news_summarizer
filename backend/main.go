package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	_ "modernc.org/sqlite"
)

// Глобальные переменные для отслеживания состояния пайплайна
var (
	pipelineRunning = false
	pipelineStatus  = "ready" // ready, running, completed, error
	pipelineStep    = ""
	pipelineMutex   sync.RWMutex
	dataReady       = false
)

// Структуры данных
type NewsItem struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Summary     string     `json:"summary"`
	Source      string     `json:"source"`
	URL         string     `json:"url"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
	Category    string     `json:"category"`
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

// API для получения статуса пайплайна
func getPipelineStatus(w http.ResponseWriter, r *http.Request) {
	pipelineMutex.RLock()
	status := map[string]interface{}{
		"running":   pipelineRunning,
		"status":    pipelineStatus,
		"step":      pipelineStep,
		"dataReady": dataReady,
	}
	pipelineMutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// Автоматический пайплайн обработки новостей
func runNewsPipeline() {
	pipelineMutex.Lock()
	pipelineRunning = true
	pipelineStatus = "running"
	dataReady = false
	pipelineMutex.Unlock()

	log.Printf("=== Начало пайплайна обработки новостей ===")
	
	// Получаем корневую директорию проекта
	workingDir, _ := os.Getwd()
	var projectRoot string
	if strings.Contains(workingDir, "backend") {
		projectRoot = filepath.Dir(workingDir)
	} else {
		projectRoot = workingDir
	}
	
	// Шаг 1: Telegram парсер
	setPipelineStep("Запуск Telegram парсера...")
	log.Printf("Шаг 1: Запуск Telegram парсера...")
	
	tgParserDir := filepath.Join(projectRoot, "parser", "tg_parser")
	cmd := exec.Command("python", "parser_prototype.py")
	cmd.Dir = tgParserDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Ошибка TG парсера: %v, output: %s", err, string(output))
	} else {
		log.Printf("TG парсер выполнен успешно")
	}
	
	// Шаг 2: Website парсер
	setPipelineStep("Запуск Website парсера...")
	log.Printf("Шаг 2: Запуск Website парсера...")
	
	websiteParserDir := filepath.Join(projectRoot, "parser", "website_parser")
	if _, err := os.Stat(filepath.Join(websiteParserDir, "script.py")); err == nil {
		cmd = exec.Command("python", "script.py")
		cmd.Dir = websiteParserDir
		output, err = cmd.CombinedOutput()
		if err != nil {
			log.Printf("Ошибка Website парсера: %v, output: %s", err, string(output))
		} else {
			log.Printf("Website парсер выполнен успешно")
		}
	} else {
		log.Printf("Website парсер не найден, пропускаем")
	}
	
	// Шаг 3: ML обработка
	setPipelineStep("ML обработка и классификация...")
	log.Printf("Шаг 3: Запуск ML обработки...")
	
	mlScriptDir := filepath.Join(projectRoot, "ml", "scripts")
	cmd = exec.Command("python", "main.py")
	cmd.Dir = mlScriptDir
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Printf("Ошибка ML скрипта: %v, output: %s", err, string(output))
		setPipelineStep("Ошибка ML обработки")
		pipelineMutex.Lock()
		pipelineStatus = "error"
		pipelineRunning = false
		pipelineMutex.Unlock()
		return
	} else {
		log.Printf("ML скрипт выполнен успешно")
	}
	
	// Шаг 4: Проверка готовности данных
	setPipelineStep("Подготовка данных для отображения...")
	log.Printf("Шаг 4: Проверка готовности данных...")
	
	filteredDataPath := filepath.Join(projectRoot, "ml", "filtered_data.json")
	if _, err := os.Stat(filteredDataPath); err == nil {
		log.Printf("Данные готовы для отображения")
		
		pipelineMutex.Lock()
		pipelineStatus = "completed"
		pipelineStep = "Готово! Данные загружены"
		pipelineRunning = false
		dataReady = true
		pipelineMutex.Unlock()
	} else {
		log.Printf("Файл filtered_data.json не найден")
		setPipelineStep("Ошибка: данные не готовы")
		pipelineMutex.Lock()
		pipelineStatus = "error"
		pipelineRunning = false
		pipelineMutex.Unlock()
	}
	
	log.Printf("=== Пайплайн обработки новостей завершен ===")
}

// Вспомогательная функция для обновления статуса
func setPipelineStep(step string) {
	pipelineMutex.Lock()
	pipelineStep = step
	pipelineMutex.Unlock()
}

func main() {
	// Инициализация БД
	initDB()

	// Запускаем пайплайн в фоне
	log.Println("Запускаем автоматический пайплайн обработки новостей...")
	go runNewsPipeline()

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

	// Статус пайплайна
	api.HandleFunc("/pipeline/status", getPipelineStatus).Methods("GET")

	// ML эндпоинты
	api.HandleFunc("/summarize", summarizeText).Methods("POST")
	// Запуск парсеров/обновления БД (может выполняться асинхронно)
	api.HandleFunc("/refresh", refreshHandler).Methods("POST")

	// Статические файлы фронтенда
	frontendDir := filepath.Join("..", "frontend")
	r.PathPrefix("/frontend/").Handler(http.StripPrefix("/frontend/", http.FileServer(http.Dir(frontendDir))))

	// Redirect с корня на frontend
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/frontend/index.html", http.StatusMovedPermanently)
	})

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
}

// Инициализация базы данных
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
	// Настройки для устойчивой параллельной работы
	_, _ = db.Exec("PRAGMA journal_mode=WAL;")
	_, _ = db.Exec("PRAGMA busy_timeout=5000;")
	db.SetMaxOpenConns(1)
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
	
	// Проверяем готовность данных
	pipelineMutex.RLock()
	isDataReady := dataReady
	currentStatus := pipelineStatus
	pipelineMutex.RUnlock()
	
	if !isDataReady {
		// Если данные не готовы, возвращаем статус пайплайна
		response := map[string]interface{}{
			"status": "loading",
			"pipeline_status": currentStatus,
			"message": "Данные обрабатываются, пожалуйста подождите",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// Данные готовы, читаем из filtered_data.json
	workingDir, _ := os.Getwd()
	var projectRoot string
	if strings.Contains(workingDir, "backend") {
		projectRoot = filepath.Dir(workingDir)
	} else {
		projectRoot = workingDir
	}
	
	filteredDataPath := filepath.Join(projectRoot, "ml", "filtered_data.json")
	
	// Читаем файл
	data, err := os.ReadFile(filteredDataPath)
	if err != nil {
		log.Printf("Ошибка чтения filtered_data.json: %v", err)
		http.Error(w, "Данные не найдены", http.StatusNotFound)
		return
	}
	
	// Парсим JSON структуру filtered_data.json
	type FilteredMessage struct {
		Text       string  `json:"text"`
		Date       string  `json:"date"`
		Link       string  `json:"link"`
		Category   string  `json:"category"`
		Confidence float64 `json:"confidence"`
	}
	
	type FilteredChannel struct {
		ChannelName string            `json:"channel_name"`
		Messages    []FilteredMessage `json:"messages"`
	}
	
	var filteredData []FilteredChannel
	if err := json.Unmarshal(data, &filteredData); err != nil {
		log.Printf("Ошибка парсинга filtered_data.json: %v", err)
		http.Error(w, "Ошибка обработки данных", http.StatusInternalServerError)
		return
	}
	
	// Получаем параметры запроса
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "20"
	}
	limitInt, _ := strconv.Atoi(limit)
	
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	pageInt, _ := strconv.Atoi(page)
	if pageInt < 1 {
		pageInt = 1
	}
	
	categoriesParam := r.URL.Query().Get("categories")
	var filterCategories []string
	if categoriesParam != "" {
		filterCategories = strings.Split(categoriesParam, ",")
		for i := range filterCategories {
			filterCategories[i] = strings.TrimSpace(filterCategories[i])
		}
	}
	
	// Преобразуем в формат NewsItem
	var allNews []NewsItem
	var filteredNews []NewsItem
	id := 1
	
	for _, channel := range filteredData {
		for _, message := range channel.Messages {
			// Создаем NewsItem
			item := NewsItem{
				ID:       id,
				Title:    fmt.Sprintf("Новость от %s", channel.ChannelName),
				Content:  message.Text,
				Summary:  message.Text,
				Source:   channel.ChannelName,
				URL:      message.Link,
				Category: message.Category,
			}
			
			// Пытаемся парсить дату
			if message.Date != "" {
				if parsed, err := time.Parse("15:04:05", message.Date); err == nil {
					// Если только время, добавляем сегодняшнюю дату
					today := time.Now().Format("2006-01-02")
					fullTime := today + "T" + parsed.Format("15:04:05") + "Z"
					if parsedFull, err := time.Parse(time.RFC3339, fullTime); err == nil {
						item.PublishedAt = &parsedFull
						item.CreatedAt = parsedFull
					}
				}
			}
			
			allNews = append(allNews, item)
			id++
		}
	}
	
	// Фильтруем по категориям если указаны
	if len(filterCategories) > 0 {
		for _, item := range allNews {
			categoryMatch := false
			for _, cat := range filterCategories {
				if strings.EqualFold(item.Category, cat) {
					categoryMatch = true
					break
				}
			}
			if categoryMatch {
				filteredNews = append(filteredNews, item)
			}
		}
	} else {
		filteredNews = allNews
	}
	
	// Рассчитываем пагинацию
	totalItems := len(filteredNews)
	totalPages := (totalItems + limitInt - 1) / limitInt
	startIndex := (pageInt - 1) * limitInt
	endIndex := startIndex + limitInt
	
	if startIndex >= totalItems {
		startIndex = totalItems
	}
	if endIndex > totalItems {
		endIndex = totalItems
	}
	
	// Получаем новости для текущей страницы
	var paginatedNews []NewsItem
	if startIndex < endIndex {
		paginatedNews = filteredNews[startIndex:endIndex]
	}
	
	// Формируем ответ с информацией о пагинации
	response := map[string]interface{}{
		"news": paginatedNews,
		"pagination": map[string]interface{}{
			"currentPage": pageInt,
			"totalPages":  totalPages,
			"totalItems":  totalItems,
			"limit":       limitInt,
			"hasNext":     pageInt < totalPages,
			"hasPrev":     pageInt > 1,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Обработчик /api/refresh — запускает парсеры и обновляет БД
func refreshHandler(w http.ResponseWriter, r *http.Request) {
	// Опция sync=true позволит дождаться завершения; иначе запускаем асинхронно
	syncMode := r.URL.Query().Get("sync") == "true"
	sentences := 2
	if s := r.URL.Query().Get("sentences"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			sentences = v
		}
	}

	// categories передаем в job через query
	categories := r.URL.Query().Get("categories")

	if syncMode {
		// Выполняем синхронно
		runRefreshJob(categories, sentences)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "completed"})
		return
	}

	// Фоново запускаем задачу
	go func(cats string, sent int) {
		runRefreshJob(cats, sent)
	}(categories, sentences)

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})
}

// Основная задача обновления: запускаем парсеры, вставляем новости и запускаем backfill суммаризаций
func runRefreshJob(categories string, sentences int) {
	log.Printf("Запущена задача refresh (categories=%s, sentences=%d)", categories, sentences)

	// Запускаем парсеры (если есть скрипты в папке parser)
	parsers := []string{"../parser/tg_parser/parser_prototype.py", "../parser/website_parser/script.py"}
	for _, p := range parsers {
		cmd := exec.Command("python", p)
		cmd.Dir = ""
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Парсер %s завершился с ошибкой: %v, output: %s", p, err, string(out))
			continue
		}
		log.Printf("Парсер %s выполнен, output: %s", p, string(out))
	}

	// После парсеров — запускаем чтение их результатов и вставку в БД.
	// Упростим: парсеры записывают JSON в known locations — website_parser -> data/all_information.json, tg_parser -> mocks/export.json
	// website
	websiteFile := "parser/website_parser/data/all_information.json"
	if b, err := os.ReadFile(websiteFile); err == nil {
		var arr []map[string]interface{}
		if err := json.Unmarshal(b, &arr); err == nil {
			for _, it := range arr {
				title, _ := it["title"].(string)
				content, _ := it["content"].(string)
				url, _ := it["url"].(string)
				category, _ := it["category"].(string)
				source := "web"
				// Генерируем summary сразу (блокируемый, но короткий)
				summary, err := callPythonSummarizerWithTimeout(content, sentences, 10*time.Second)
				if err != nil {
					log.Printf("Ошибка суммаризации web %s: %v", url, err)
				}
				query := `INSERT INTO news (title, content, summary, source, url, published_at, category) VALUES (?, ?, ?, ?, ?, ?, ?)`
				_, err = db.Exec(query, title, content, summary, source, url, nil, category)
				if err != nil {
					log.Printf("Не добавлена новость %s: %v", url, err)
				} else {
					log.Printf("Вставляем в БД web: %s", url)
				}
			}
		}
	}

	// tg
	tgFile := "parser/tg_parser/mocks/export.json"
	if b, err := os.ReadFile(tgFile); err == nil {
		var arr []map[string]interface{}
		if err := json.Unmarshal(b, &arr); err == nil {
			for _, it := range arr {
				title, _ := it["title"].(string)
				content, _ := it["content"].(string)
				url, _ := it["url"].(string)
				channel, _ := it["channel"].(string)
				source := channel
				summary, err := callPythonSummarizerWithTimeout(content, sentences, 10*time.Second)
				if err != nil {
					log.Printf("Ошибка суммаризации tg %s: %v", url, err)
				}
				query := `INSERT INTO news (title, content, summary, source, url, published_at, category) VALUES (?, ?, ?, ?, ?, ?, ?)`
				_, err = db.Exec(query, title, content, summary, source, url, nil, "ТГ")
				if err != nil {
					log.Printf("Не добавлена новость (возможно дубликат) %s: %v", url, err)
				} else {
					log.Printf("Вставляем в БД TG: %s", url)
				}
			}
		}
	}

	// Фоновая добивка summary для записей без summary
	go func() {
		log.Println("Запускаем backfill summary для существующих записей...")
		rows, err := db.Query("SELECT id, content FROM news WHERE summary IS NULL OR summary = ''")
		if err != nil {
			log.Printf("backfill query error: %v", err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var id int
			var content string
			if err := rows.Scan(&id, &content); err != nil {
				continue
			}
			log.Printf("backfill: генерируем summary для ID=%d...", id)
			s, err := callPythonSummarizerWithTimeout(content, sentences, 10*time.Second)
			if err != nil {
				log.Printf("backfill summarize error for id=%d: %v", id, err)
				continue
			}
			_, err = db.Exec("UPDATE news SET summary = ? WHERE id = ?", s, id)
			if err != nil {
				log.Printf("backfill update error for id=%d: %v", id, err)
			}
		}
	}()

	log.Println("Задача refresh завершена")
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
	// Опционально можно задать количество предложений через query param sentences (default 2)
	sentences := 2
	if s := r.URL.Query().Get("sentences"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			sentences = v
		}
	}

	// Вызов Python скрипта для суммаризации с таймаутом
	summary, err := callPythonSummarizerWithTimeout(req.Content, sentences, 15*time.Second)
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
func callPythonSummarizer(content string, sentenceCount int, pythonPath string) (string, error) {
	// Подготавливаем входные данные для Python скрипта
	input := map[string]interface{}{
		"content":        content,
		"sentence_count": sentenceCount,
	}

	inputJSON, err := json.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("ошибка подготовки данных: %v", err)
	}

	mlScript := os.Getenv("ML_SCRIPT_PATH")
	if mlScript == "" {
		mlScript = filepath.Join("..", "ml", "scripts", "main.py")
	}

	// Позволим указать python интерпретатор, иначе попробуем 'python'
	pythonBin := pythonPath
	if pythonBin == "" {
		pythonBin = "python"
	}

	cmd := exec.Command(pythonBin, mlScript)
	cmd.Stdin = bytes.NewReader(inputJSON)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения Python скрипта: %v, output: %s", err, string(out))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(out, &result); err != nil {
		return "", fmt.Errorf("ошибка парсинга результата: %v, output: %s", err, string(out))
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

// Вызов суммаризатора с контекстом/таймаутом
func callPythonSummarizerWithTimeout(content string, sentenceCount int, timeout time.Duration) (string, error) {
	// Попытаемся использовать ML_PYTHON из окружения (полезно при venv)
	pythonPath := os.Getenv("ML_PYTHON")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	type resultT struct {
		s   string
		err error
	}

	ch := make(chan resultT, 1)
	go func() {
		s, e := callPythonSummarizer(content, sentenceCount, pythonPath)
		ch <- resultT{s, e}
	}()

	select {
	case <-ctx.Done():
		return "", fmt.Errorf("summarizer timeout after %v", timeout)
	case res := <-ch:
		return res.s, res.err
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
