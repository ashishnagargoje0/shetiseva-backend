package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var resetSecretKey = []byte("your-secret-reset-key") // 🔐 Change in production

// ✅ Custom claims for reset token
type ResetClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// ✅ GenerateResetToken creates a JWT for password reset
func GenerateResetToken(userID string) (string, error) {
	claims := &ResetClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(resetSecretKey)
}

// ✅ ValidateResetToken checks if the token is valid and returns userID
func ValidateResetToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ResetClaims{}, func(token *jwt.Token) (interface{}, error) {
		return resetSecretKey, nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(*ResetClaims)
	if !ok || claims.UserID == "" {
		return "", errors.New("invalid reset token claims")
	}

	return claims.UserID, nil
}
