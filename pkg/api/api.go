package api

import (
	"net/http"

	"github.com/Perehodko/final_project_go/pkg/middleware"
)

// Init инициализирует обработчики API
func Init() {
	// Обработчик аутентификации (без middleware)
	http.HandleFunc("/api/signin", SignInHandler)

	// Защищенные обработчики с middleware
	http.HandleFunc("/api/nextdate", middleware.AuthMiddleware(NextDateHandler))
	http.HandleFunc("/api/task", middleware.AuthMiddleware(taskHandler))
	http.HandleFunc("/api/tasks", middleware.AuthMiddleware(tasksHandler))
	http.HandleFunc("/api/task/done", middleware.AuthMiddleware(doneHandler))
}
