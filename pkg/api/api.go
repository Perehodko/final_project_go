package api

import "net/http"

// Init регистрирует все обработчики API
func Init() {
	http.HandleFunc("/api/nextdate", NextDateHandler)
	// Здесь будут регистрироваться другие обработчики по мере реализации
	// http.HandleFunc("/api/tasks", TasksHandler)
	// http.HandleFunc("/api/task", TaskHandler)
	// http.HandleFunc("/api/task/done", TaskDoneHandler)
	// http.HandleFunc("/api/login", LoginHandler)
}
