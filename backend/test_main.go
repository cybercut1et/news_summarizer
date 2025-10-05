package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	_ "modernc.org/sqlite"
)

type NewsItem struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var db *sql.DB

func initDB() {
	var err error
	log.Println("Инициализация базы данных...")
	db, err = sql.Open("sqlite", "./news.db")
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка проверки соединения с БД:", err)
	}
	log.Println("База данных готова к работе")

	// Создаем таблицу если не существует
	createTable := `CREATE TABLE IF NOT EXISTS news (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT,
		summary TEXT,
		source TEXT,
		url TEXT,
		published_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		category TEXT DEFAULT 'Общее'
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal("Ошибка создания таблицы:", err)
	}
}

func getNews(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос на новости")

	rows, err := db.Query("SELECT id, title, content FROM news LIMIT 10")
	if err != nil {
		log.Printf("Ошибка запроса: %v", err)
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var news []NewsItem
	for rows.Next() {
		var item NewsItem
		err := rows.Scan(&item.ID, &item.Title, &item.Content)
		if err != nil {
			log.Printf("Ошибка сканирования: %v", err)
			continue
		}
		news = append(news, item)
	}

	log.Printf("Найдено новостей: %d", len(news))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}

func main() {
	initDB()

	r := mux.NewRouter()

	// Простые эндпоинты для тестирования
	r.HandleFunc("/health", healthCheck).Methods("GET")
	r.HandleFunc("/api/news", getNews).Methods("GET")

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

	fmt.Printf("Тестовый сервер запущен на порту %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
