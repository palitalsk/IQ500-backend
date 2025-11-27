package utils

import (
	"errors"
	"fmt"
	"main/domain"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type authCustomClaimsMember struct {
	ID       string `bson:"_id" json:"id"`
	Username string `bson:"username" json:"username"`
	Active   bool   `bson:"active" json:"active"`

	jwt.StandardClaims
}

func Authentication(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}
		authHeader(c, db)
		c.Next()
	}
}

func authHeader(c *gin.Context, db *mongo.Database) {
	bearToken, errToken := extractToken(c)
	if errToken != nil {
		fmt.Println("authorization header is not valid[1]", errToken)
		c.AbortWithStatusJSON(401, gin.H{
			"code":  401,
			"error": "ไม่สามารถใช้งานได้",
		})
		return
	}

	if len(bearToken) != 2 {
		fmt.Println("authorization header is not valid[2]")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "error": "ไม่สามารถใช้งานได้"})
		return
	}

	tokenString := bearToken[1]
	var claims = &authCustomClaimsMember{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	if err != nil {
		fmt.Println("authorization header is not valid[3]", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "error": "ไม่สามารถใช้งานได้"})
		return
	}

	if claims.Username == "" {
		fmt.Println("authorization header is not valid[4]")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "error": "ไม่สามารถใช้งานได้"})
		return
	}

	id, _ := primitive.ObjectIDFromHex(claims.ID)
	var user domain.User

	option := options.FindOne()
	option.SetSort(nil)

	filter := bson.M{
		"_id": id,
	}

	err = db.Collection("user_account").FindOne(c, filter, option).Decode(&user)
	if err != nil {
		fmt.Println("authorization header is not valid[5]", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": 401, "error": "ไม่พบข้อมูลผู้ใช้งาน"})
		return
	}

	if tokenString != user.Token {
		fmt.Println("authorization header is not valid[6]")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": 401, "error": "token หมดอายุ"})
		return
	}

	if !user.Active {
		fmt.Println("authorization header is not valid[7]")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": 401, "error": "บัญชีของคุณถูกระงับการใช้งาน"})
		return
	}

	c.Set("user_id", user.ID.Hex())
}

func extractToken(c *gin.Context) ([]string, error) {
	var token []string
	bearToken := c.Request.Header.Get("Authorization")
	if !strings.Contains(bearToken, "Bearer") {
		return token, errors.New("ERROR BEARER TOKEN")
	}
	token = strings.Split(bearToken, "Bearer ")
	return token, nil
}
