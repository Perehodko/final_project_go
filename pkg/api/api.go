package api

import "net/http"

// Init регистрирует все обработчики API
func Init() {
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	// http.HandleFunc("/api/tasks", TasksHandler)
	// http.HandleFunc("/api/task", TaskHandler)
	// http.HandleFunc("/api/task/done", TaskDoneHandler)
	// http.HandleFunc("/api/login", LoginHandler)
}

// taskHandler обрабатывает все запросы к /api/task
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	default:
		writeJSONError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
