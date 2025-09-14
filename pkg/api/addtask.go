package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Perehodko/final_project_go/pkg/db"
)

// addTaskHandler обрабатывает POST запросы для добавления задач
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Парсим JSON из тела запроса
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, "Ошибка разбора JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Валидируем и обрабатываем задачу
	if err := processTask(&task); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Добавляем задачу в БД
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ (конвертируем int64 в string)
	writeJSONResponse(w, map[string]interface{}{"id": strconv.FormatInt(id, 10)}, http.StatusOK)
}

// processTask обрабатывает и валидирует данные задачи
func processTask(task *db.Task) error {
	// Проверяем обязательное поле title
	if task.Title == "" {
		return fmt.Errorf("не указан заголовок задачи")
	}

	// СНАЧАЛА проверяем правило повторения (если указано)
	if task.Repeat != "" {
		// Проверяем правило повторения с помощью NextDate
		// Используем текущую дату как базовую для проверки
		now := time.Now()
		currentDate := now.Format("20060102")
		_, err := NextDate(now, currentDate, task.Repeat)
		if err != nil {
			return fmt.Errorf("неверный формат правила повторения: %v", err)
		}
	}

	now := time.Now()
	currentDate := now.Format("20060102")

	// Если дата не указана, используем текущую
	if task.Date == "" {
		task.Date = currentDate
		return nil
	}

	// Проверяем формат даты
	parsedDate, err := time.Parse("20060102", task.Date)
	if err != nil {
		return fmt.Errorf("неверный формат даты")
	}

	// Если дата в прошлом
	if !afterNow(parsedDate, now) {
		// ВСЕГДА используем текущую дату, независимо от правила повторения
		task.Date = currentDate
	}

	return nil
}

// writeJSONResponse записывает JSON ответ
func writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Ошибка кодирования JSON", http.StatusInternalServerError)
	}
}

// writeJSONError записывает JSON ошибку
func writeJSONError(w http.ResponseWriter, errorMsg string, statusCode int) {
	writeJSONResponse(w, map[string]string{"error": errorMsg}, statusCode)
}
