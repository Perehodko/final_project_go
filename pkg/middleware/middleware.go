package middleware

import (
	"net/http"
	"os"

	"github.com/Perehodko/final_project_go/pkg/auth"
)

// AuthMiddleware middleware для проверки аутентификации
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwtToken string

			// Получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				jwtToken = cookie.Value
			}

			// Если нет куки, проверяем Authorization header (для тестов)
			if jwtToken == "" {
				authHeader := r.Header.Get("Authorization")
				if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
					jwtToken = authHeader[7:]
				}
			}

			var valid bool

			// Валидируем JWT-токен
			if jwtToken != "" {
				valid, err = auth.ValidateToken(jwtToken)
				if err != nil {
					valid = false
				}
			}

			if !valid {
				// Возвращаем ошибку авторизации 401
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		}

		// Продолжаем выполнение к основному обработчику
		next(w, r)
	})
}
