package jwt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	gjwt "github.com/golang-jwt/jwt/v5"
)

var secret = []byte("2345-qwer-0987-lkjh7755jnfnvfnv")

type Claims struct {
	UserID int    `json:"userID"`
	Email  string `json:"email"`
	Role   string `json:"role"`

	gjwt.RegisteredClaims
}

func GenerateToken(userID int, email, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: gjwt.RegisteredClaims{
			ExpiresAt: gjwt.NewNumericDate(
				time.Now().Add(15 * time.Minute)),
		},
	}

	token := gjwt.NewWithClaims(gjwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := gjwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *gjwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*gjwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
