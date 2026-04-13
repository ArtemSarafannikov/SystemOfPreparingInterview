package auth_token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

var secretKey = []byte("test-secret-key")

type JWTClaims struct {
	UserID    uuid.UUID
	ExpiredAt time.Time
}

func (j JWTClaims) GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": j.UserID,
		"exp":     j.ExpiredAt.Unix(),
	})

	return token.SignedString(secretKey)
}

// GetClaimsFromJWTToken Проверяет токен на корректность, валидность и просроченность
//
// Если все корректно, то возвращает структуру JWTClaims с данными из токена
func GetClaimsFromJWTToken(tokenString string) (JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return JWTClaims{}, fmt.Errorf("jwt.Parse: %w", err)
	}

	if !token.Valid {
		return JWTClaims{}, fmt.Errorf("invalid token")
	}

	switch claims := token.Claims.(type) {
	case jwt.MapClaims:
		return buildJWTClaims(claims)
	default:
		return JWTClaims{}, fmt.Errorf("invalid claims")
	}
}

func buildJWTClaims(claims jwt.MapClaims) (JWTClaims, error) {
	expiredAt, err := claims.GetExpirationTime()
	if err != nil {
		return JWTClaims{}, fmt.Errorf("claims.GetExpirationTime: %w", err)
	}

	return JWTClaims{
		UserID:    uuid.MustParse(claims["user_id"].(string)),
		ExpiredAt: expiredAt.Time,
	}, nil
}
