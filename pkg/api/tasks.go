package api

import (
	"encoding/json"
	"net/http"

	"github.com/Perehodko/final_project_go/pkg/db"
)

// TasksResp структура для ответа с задачами
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler обрабатывает GET запросы для получения списка задач
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметр поиска
	search := r.URL.Query().Get("search")

	var tasks []*db.Task
	var err error

	if search != "" {
		// Если есть параметр поиска, используем поисковую функцию
		tasks, err = db.SearchTasks(search, 50)
	} else {
		// Иначе получаем все задачи с лимитом
		tasks, err = db.Tasks(50)
	}

	if err != nil {
		writeJSONError(w, "Ошибка получения задач: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Если tasks nil, создаем пустой слайс для корректного JSON
	if tasks == nil {
		tasks = []*db.Task{}
	}

	writeJSONResponse(w, TasksResp{Tasks: tasks}, http.StatusOK)
}

// taskHandler обрабатывает все запросы к /api/task
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r) 
	default:
		writeJSONError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// getTaskHandler обрабатывает GET запросы для получения задачи по ID
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSONError(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	writeJSONResponse(w, task, http.StatusOK)
}

// updateTaskHandler обрабатывает PUT запросы для обновления задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Парсим JSON из тела запроса
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, "Ошибка разбора JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем наличие ID
	if task.ID == "" {
		writeJSONError(w, "Не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	// Валидируем и обрабатываем задачу
	if err := processTask(&task); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обновляем задачу в БД
	err := db.UpdateTask(&task)
	if err != nil {
		writeJSONError(w, "Ошибка обновления задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	writeJSONResponse(w, map[string]interface{}{}, http.StatusOK)
}

// deleteTaskHandler обрабатывает DELETE запросы для удаления задачи
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeJSONError(w, "Ошибка удаления задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	writeJSONResponse(w, map[string]interface{}{}, http.StatusOK)
}
