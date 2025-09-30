package utils

import (
	"time"

	"github.com/febryanhernanda/social-media-apps/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret []byte
}

func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{
		secret: []byte(secret),
	}
}

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

/* Generate Token */
func (j *JWTManager) GenerateToken(user *models.User) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "sosmed",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

/* Validate token */
func (j *JWTManager) ValidateToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
