package gateway

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AccessClamis struct {
	UserID   uint
	Role     string
	Username string
	jwt.RegisteredClaims
}

func ValidateAccessToken(tokenString string) (*AccessClamis, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&AccessClamis{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("Wrong token method")
			}
			return LoadSecretKey(), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessClamis)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.ExpiresAt == nil || claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}

func AuthGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":     401,
				"messeage": "Unauthoried",
				"data":     nil,
			})
			return
		}

		token := strings.TrimPrefix(authorization, "Bearer ")
		if token == authorization {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":     401,
				"messeage": "Token wrong format",
				"data":     nil,
			})
			return
		}

		claims, err := ValidateAccessToken(token)
		if err != nil {
			log.Println("token invalid:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":     401,
				"messeage": err.Error(),
				"data":     nil,
			})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}
