package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

// InitJWT initializes JWT secret key
func InitJWT(secret string) {
	jwtSecret = []byte(secret)
}

// Claims represents JWT claims
type Claims struct {
	Sub         string   `json:"sub"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token with specified duration
func GenerateToken(userID, role string, permissions []string, duration time.Duration) (string, error) {
	if len(jwtSecret) == 0 {
		return "", errors.New("JWT secret not initialized")
	}

	now := time.Now()
	expirationTime := now.Add(duration)

	claims := &Claims{
		Sub:         userID,
		Role:        role,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyToken verifies and parses JWT token
func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	if len(jwtSecret) == 0 {
		return nil, errors.New("JWT secret not initialized")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return jwt.MapClaims{
			"sub":         claims.Sub,
			"username":    claims.Username,
			"email":       claims.Email,
			"role":        claims.Role,
			"permissions": claims.Permissions,
		}, nil
	}

	return nil, errors.New("invalid token claims")
}

// GetClaimsFromToken extracts claims from token
func GetClaimsFromToken(tokenString string) (*Claims, error) {
	if len(jwtSecret) == 0 {
		return nil, errors.New("JWT secret not initialized")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}
