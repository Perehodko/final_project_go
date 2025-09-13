package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Perehodko/final_project_go/pkg/db"
	"github.com/Perehodko/final_project_go/pkg/server"
	"github.com/Perehodko/final_project_go/tests"
)

func main() {
	// Путь к БД - в родительской директории (где тесты ожидают)
	dbFile := "scheduler.db"

	// Но если есть переменная окружения, используем её
	if envFile := os.Getenv("TODO_DBFILE"); envFile != "" {
		dbFile = envFile
	}

	// Получаем абсолютный путь для логов
	absPath, err := filepath.Abs(dbFile)
	if err != nil {
		log.Fatalf("Ошибка получения абсолютного пути: %v", err)
	}

	log.Printf("Используем файл БД: %s", absPath)

	err = db.Init(dbFile)
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer db.Close()

	// Получаем порт
	port := getPort()
	webDir := "./web"

	// Запускаем сервер
	err = server.StartServer(port, webDir)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

func getPort() int {
	// 1. Проверяем переменную окружения TODO_PORT
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		if port, err := strconv.Atoi(envPort); err == nil && port > 0 && port < 65536 {
			log.Printf("Используем порт из переменной окружения TODO_PORT: %d", port)
			return port
		} else {
			log.Printf("Некорректное значение TODO_PORT: %s, используем порт по умолчанию", envPort)
		}
	}

	// 2. Используем порт из tests/settings.go
	log.Printf("Используем порт по умолчанию: %d", tests.Port)
	return tests.Port
}
