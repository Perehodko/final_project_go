package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Perehodko/final_project_go/pkg/api"
)

func StartServer(port int, webDir string) error {
	// Инициализируем API обработчики
	api.Init()

	// Обслуживание статических файлов
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("Запуск веб-сервера на порту %d", port)
	log.Printf("Откройте http://localhost:%d/ в браузере", port)

	// Проверяем, включена ли аутентификация
	if password := os.Getenv("TODO_PASSWORD"); password != "" {
		fmt.Println("Authentication is ENABLED")
	} else {
		fmt.Println("Authentication is DISABLED")
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
