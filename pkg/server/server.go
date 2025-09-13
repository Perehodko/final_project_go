package server

import (
	"fmt"
	"log"
	"net/http"
)

func StartServer(port int, webDir string) error {
	// Обслуживание статических файлов
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	// Здесь позже будут добавлены API endpoints
	// http.HandleFunc("/api/tasks", api.GetTasksHandler)
	// http.HandleFunc("/api/tasks/add", api.AddTaskHandler)

	log.Printf("Запуск веб-сервера на порту %d", port)
	log.Printf("Откройте http://localhost:%d/ в браузере", port)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
