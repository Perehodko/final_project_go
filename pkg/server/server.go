package server

import (
    "fmt"
    "log"
    "net/http"
)

func StartServer(port int, webDir string) error {
    http.Handle("/", http.FileServer(http.Dir(webDir)))
    
    log.Printf("Запуск веб-сервера на порту %d", port)
    log.Printf("Откройте http://localhost:%d/ в браузере", port)
    
    return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
