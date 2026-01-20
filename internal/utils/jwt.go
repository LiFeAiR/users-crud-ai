package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims расширенные claims для JWT токена
type Claims struct {
	UserID      int      `json:"user_id"`
	Email       string   `json:"email"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// GenerateJWT генерирует JWT токен для пользователя
func GenerateJWT(secretKey string, userID int, email, name string, permissions []string) (string, error) {
	// Устанавливаем срок действия токена (например, 24 часа)
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID:      userID,
		Email:       email,
		Name:        name,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}

// ValidateJWT проверяет валидность JWT токена
func ValidateJWT(tokenString, secretKey string) (*Claims, error) {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !tkn.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
