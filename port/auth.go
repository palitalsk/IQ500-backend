package port

import (
	"main/domain"

	"github.com/gin-gonic/gin"
)

type AuthenticationService interface {
	Login(c *gin.Context, payload domain.PayloadUser) (string, error)
	Register(c *gin.Context, payload domain.PayloadUser) error
	ResetPassword(c *gin.Context, payload domain.PayloadResetPassword) error
}





