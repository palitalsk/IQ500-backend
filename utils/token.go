package utils

import (
	"fmt"
	"main/domain"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func GenTokenMember(c *gin.Context, body domain.User) (string, error) {
	mySigningKey := []byte(os.Getenv("ACCESS_SECRET"))
	day := time.Hour * 24

	claims := &authCustomClaimsMember{
		body.ID.Hex(),
		body.Username,
		body.Active,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(day * 366).Unix(),
			Issuer:    "system",
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Println(err)
	}

	return ss, nil
}
