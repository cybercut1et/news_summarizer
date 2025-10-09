package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Устанавливаем переменные окружения для путей
	backendDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("Не удалось получить путь к бэкенду: %v", err)
	}

	projectRoot := filepath.Dir(backendDir)

	os.Setenv("DB_PATH", filepath.Join(backendDir, "news.db"))
	os.Setenv("ML_SCRIPT_PATH", filepath.Join(projectRoot, "ml", "scripts", "main.py"))
	os.Setenv("PORT", "8081")

	// Запускаем main.go
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = backendDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("Запуск сервера на порту %s...", os.Getenv("PORT"))
	log.Printf("DB_PATH: %s", os.Getenv("DB_PATH"))
	log.Printf("ML_SCRIPT_PATH: %s", os.Getenv("ML_SCRIPT_PATH"))

	err = cmd.Run()
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
