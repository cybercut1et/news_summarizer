package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "./news.db")
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	fmt.Println("Содержимое таблицы 'news':")
	rows, err := db.Query("SELECT id, title, source, category, summary FROM news LIMIT 10")
	if err != nil {
		log.Fatal("Ошибка запроса к таблице news:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var title, source, category string
		var summary sql.NullString
		if err := rows.Scan(&id, &title, &source, &category, &summary); err != nil {
			log.Println("Ошибка сканирования строки:", err)
			continue
		}
		summaryStr := "NULL"
		if summary.Valid {
			summaryStr = summary.String
			if len(summaryStr) > 50 {
				summaryStr = summaryStr[:50] + "..."
			}
		}
		fmt.Printf("ID: %d, Title: %s, Source: %s, Category: %s, Summary: %s\n", id, title, source, category, summaryStr)
	}

	fmt.Println("\nСодержимое таблицы 'users':")
	rows, err = db.Query("SELECT id, uuid, sources FROM users LIMIT 10")
	if err != nil {
		log.Fatal("Ошибка запроса к таблице users:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var uuid, sources string
		if err := rows.Scan(&id, &uuid, &sources); err != nil {
			log.Println("Ошибка сканирования строки:", err)
			continue
		}
		fmt.Printf("ID: %d, UUID: %s, Sources: %s\n", id, uuid, sources)
	}
}
