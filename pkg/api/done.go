package api

import (
	"net/http"
	"time"

	"github.com/Perehodko/final_project_go/pkg/db"
)

// doneHandler обрабатывает POST запросы для отметки выполнения задачи
func doneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	// Получаем задачу по ID
	task, err := db.GetTask(id)
	if err != nil {
		writeJSONError(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	// Обрабатываем выполнение задачи
	if task.Repeat == "" {
		// Одноразовая задача - удаляем
		err = db.DeleteTask(id)
		if err != nil {
			writeJSONError(w, "Ошибка удаления задачи: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Периодическая задача - вычисляем следующую дату
		now := time.Now()

		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSONError(w, "Ошибка вычисления следующей даты: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Обновляем дату задачи
		err = db.UpdateTaskDate(id, nextDate)
		if err != nil {
			writeJSONError(w, "Ошибка обновления даты: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Возвращаем успешный ответ
	writeJSONResponse(w, map[string]interface{}{}, http.StatusOK)
}
