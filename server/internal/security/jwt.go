package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtGenerate struct {
	Secretkey string
}

func NewJwtGenerate(key string) *JwtGenerate {
	return &JwtGenerate{Secretkey: key}
}

type AccessClamis struct {
	UserID   uint
	Role     string
	Username string
	jwt.RegisteredClaims
}

func (jg *JwtGenerate) GenerateAccessToken(userID uint, role, username string) (int64, string, error) {
	claims := AccessClamis{
		UserID:   userID,
		Role:     role,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(15 * time.Minute),
			),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	response, err := token.SignedString([]byte(jg.Secretkey))
	ttl := time.Until(claims.ExpiresAt.Time).Seconds()
	return int64(ttl), response, err
}

func (jg *JwtGenerate) GenerateRefreshToken(userID uint, role, username string) (int64, string, error) {
	claims := AccessClamis{
		UserID:   userID,
		Role:     role,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(24 * 30 * time.Hour), // 1 thang
			),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	response, err := token.SignedString([]byte(jg.Secretkey))
	ttl := time.Until(claims.ExpiresAt.Time).Seconds()
	return int64(ttl), response, err
}
