package api

import (
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
