package db

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var DB *sqlx.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL DEFAULT '',
    comment TEXT,
    repeat VARCHAR(128)
);

CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

// Init инициализирует базу данных
func Init(dbFile string) error {
	// Проверяем существование файла БД
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	// Открываем базу данных с sqlx
	DB, err = sqlx.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка открытия БД: %v", err)
	}

	// Проверяем соединение
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД: %v", err)
	}

	// Если БД не существовала, создаем схему
	if install {
		log.Printf("Создание новой базы данных: %s", dbFile)
		_, err = DB.Exec(schema)
		if err != nil {
			return fmt.Errorf("ошибка создания схемы: %v", err)
		}
		log.Println("База данных успешно создана")
	} else {
		log.Printf("База данных уже существует: %s", dbFile)
	}

	return nil
}

// GetDB возвращает экземпляр БД (может пригодиться)
func GetDB() *sqlx.DB {
	return DB
}

// Close закрывает соединение с БД
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// AddTask добавляет задачу в БД
func AddTask(task *Task) (int64, error) {
	if err := task.Validate(); err != nil {
		return 0, err
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("ошибка добавления задачи: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения ID: %v", err)
	}

	return id, nil
}
