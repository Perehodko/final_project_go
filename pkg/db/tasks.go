package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Tasks возвращает список задач с ограничением по количеству
func Tasks(limit int) ([]*Task, error) {
	db := GetDB()
	if db == nil {
		return nil, fmt.Errorf("база данных не инициализирована")
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	return scanTasks(rows)
}

// SearchTasks ищет задачи по тексту или дате
func SearchTasks(search string, limit int) ([]*Task, error) {
	db := GetDB()
	if db == nil {
		return nil, fmt.Errorf("база данных не инициализирована")
	}

	// Проверяем, является ли поисковый запрос датой в формате DD.MM.YYYY
	if isDateSearch(search) {
		// Преобразуем дату из формата DD.MM.YYYY в YYYYMMDD
		date, err := convertSearchDate(search)
		if err != nil {
			return nil, err
		}

		query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`
		rows, err := db.Query(query, date, limit)
		if err != nil {
			return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
		}
		defer rows.Close()

		return scanTasks(rows)
	}

	// Поиск по тексту в title или comment
	searchPattern := "%" + search + "%"
	query := `SELECT id, date, title, comment, repeat FROM scheduler 
	          WHERE title LIKE ? OR comment LIKE ? 
	          ORDER BY date LIMIT ?`
	rows, err := db.Query(query, searchPattern, searchPattern, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	return scanTasks(rows)
}

// scanTasks сканирует строки результата запроса и возвращает список задач
func scanTasks(rows *sql.Rows) ([]*Task, error) {
	var tasks []*Task

	for rows.Next() {
		var task Task
		var id int64
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %v", err)
		}
		// Конвертируем int64 ID в string для JSON
		task.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %v", err)
	}

	return tasks, nil
}

// isDateSearch проверяет, является ли строка датой в формате DD.MM.YYYY
func isDateSearch(s string) bool {
	if len(s) != 10 {
		return false
	}
	if s[2] != '.' || s[5] != '.' {
		return false
	}

	// Проверяем, что все остальные символы - цифры
	for i, char := range s {
		if i != 2 && i != 5 && (char < '0' || char > '9') {
			return false
		}
	}

	return true
}

// convertSearchDate преобразует дату из формата DD.MM.YYYY в YYYYMMDD
func convertSearchDate(dateStr string) (string, error) {
	parts := strings.Split(dateStr, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("неверный формат даты")
	}

	day, month, year := parts[0], parts[1], parts[2]

	// Проверяем корректность даты
	_, err := time.Parse("20060102", year+month+day)
	if err != nil {
		return "", fmt.Errorf("неверная дата")
	}

	return year + month + day, nil
}
