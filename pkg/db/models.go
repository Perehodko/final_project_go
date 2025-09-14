package db

import (
	"fmt"
	"time"
)

// Task представляет задачу в системе
type Task struct {
	ID      int64  `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// TaskResponse представляет ответ API для задач
type TaskResponse struct {
	ID    int64  `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

// ToMap преобразует задачу в map для удобства работы с БД
func (t *Task) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"date":    t.Date,
		"title":   t.Title,
		"comment": t.Comment,
		"repeat":  t.Repeat,
	}
}

// Validate проверяет корректность данных задачи
func (t *Task) Validate() error {
	if t.Title == "" {
		return fmt.Errorf("не указан заголовок задачи")
	}

	// Проверяем формат даты, если она указана
	if t.Date != "" {
		_, err := time.Parse("20060102", t.Date)
		if err != nil {
			return fmt.Errorf("неверный формат даты")
		}
	}

	return nil
}
