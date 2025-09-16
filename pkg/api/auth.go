package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/Perehodko/final_project_go/pkg/auth"
)

// SignInRequest структура для запроса аутентификации
type SignInRequest struct {
	Password string `json:"password"`
}

// SignInResponse структура для ответа аутентификации
type SignInResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

// SignInHandler обработчик для аутентификации
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Получаем пароль из переменных окружения
	envPassword := os.Getenv("TODO_PASSWORD")

	// Если пароль не установлен, пропускаем аутентификацию
	if envPassword == "" {
		response := SignInResponse{Token: "no-auth"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Проверяем пароль
	if req.Password != envPassword {
		response := SignInResponse{Error: "Неверный пароль"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Генерируем JWT токен
	token, err := auth.GenerateToken()
	if err != nil {
		response := SignInResponse{Error: "Ошибка генерации токена"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := SignInResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
