package db

import "os"

// Task представляет модель задачи
type Task struct {
	ID      int64  `db:"id" json:"id"`
	Date    string `db:"date" json:"date"`       // формат YYYYMMDD
	Title   string `db:"title" json:"title"`     // VARCHAR(255)
	Comment string `db:"comment" json:"comment"` // TEXT
	Repeat  string `db:"repeat" json:"repeat"`   // VARCHAR(128)
}

// GetDBPath возвращает путь к файлу БД
func GetDBPath() string {
	// Проверяем переменную окружения TODO_DBFILE
	if dbFile := os.Getenv("TODO_DBFILE"); dbFile != "" {
		return dbFile
	}
	// Используем значение по умолчанию (совпадает с тестом)
	return "../scheduler.db"
}
