package utils

import (
	"errors"
	"go-datalaris/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte(config.JWTSecret)

func GetKey() []byte {
	return jwtKey
}

func GenerateToken(userID uint, tenantID *uint, roles []string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"tenant_id": tenantID,
		"roles":     roles,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
		"iat":       time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ===================================================================
// ✅ ParseJWT (baru ditambahkan)
// ===================================================================
// Fungsi ini digunakan untuk mem-parse token JWT dan mengambil claim-nya.
func ParseJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	return claims, nil
}

// ===================================================================
// ✅ GetClaimValue (ambil claim dari context, misalnya tenant_id, user_id)
// ===================================================================
func GetClaimValue[T any](c *gin.Context, key string) (T, error) {
	var zero T

	claimsValue, exists := c.Get("claims")
	if !exists {
		return zero, errors.New("claims not found in context")
	}

	claims, ok := claimsValue.(jwt.MapClaims)
	if !ok {
		return zero, errors.New("invalid claims type")
	}

	val, exists := claims[key]
	if !exists {
		return zero, errors.New("claim key not found: " + key)
	}

	switch v := any(val).(type) {
	case float64:
		if any(zero) == uint(0) {
			return any(uint(v)).(T), nil
		}
		if any(zero) == int(0) {
			return any(int(v)).(T), nil
		}
	}

	if casted, ok := val.(T); ok {
		return casted, nil
	}

	return zero, errors.New("claim type mismatch for key: " + key)
}
