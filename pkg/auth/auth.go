package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// CustomClaims кастомные claims для JWT
type CustomClaims struct {
	PasswordHash string `json:"pwd_hash"`
	jwt.RegisteredClaims
}

// GenerateToken генерирует JWT токен
func GenerateToken() (string, error) {
	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		return "", fmt.Errorf("password not set")
	}

	// Создаем хэш пароля для включения в токен
	hash := sha256.Sum256([]byte(password))
	passwordHash := hex.EncodeToString(hash[:])

	claims := CustomClaims{
		PasswordHash: passwordHash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Используем пароль как секрет для подписи
	return token.SignedString([]byte(password))
}

// ValidateToken проверяет валидность JWT токена
func ValidateToken(tokenString string) (bool, error) {
	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		return true, nil // Аутентификация отключена
	}

	// Парсим токен
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(password), nil
	})

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		// Проверяем, что хэш пароля в токене совпадает с текущим
		currentHash := sha256.Sum256([]byte(password))
		currentHashStr := hex.EncodeToString(currentHash[:])

		return claims.PasswordHash == currentHashStr, nil
	}

	return false, fmt.Errorf("invalid token")
}
