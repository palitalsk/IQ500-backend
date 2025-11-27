package services

import (
	"errors"
	"fmt"
	"main/domain"
	"main/port"
	"main/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthenticationService ให้บริการเกี่ยวกับ auth เดิมของ IQ500
type AuthenticationService struct {
	userRepo port.UserRepository
}

func NewAuthenticationService(
	userRepo port.UserRepository,
) *AuthenticationService {
	return &AuthenticationService{
		userRepo,
	}
}

func (s *AuthenticationService) Login(c *gin.Context, payload domain.PayloadUser) (string, error) {
	user, err := s.userRepo.GetUserByName(payload.Username)
	if err != nil {
		fmt.Println("error get user", err)
		utils.Response(c, http.StatusInternalServerError, 500, "ไม่พบยูสเซอร์ กรุณาลองใหม่อีกครั้ง", err.Error(), nil)
		return "", err
	}

	if utils.MD5Hash(payload.Password) != user.Password {
		fmt.Println("error password not match", err)
		utils.Response(c, http.StatusBadRequest, 400, "รหัสผ่านไม่ถูกต้อง กรุณาลองใหม่อีกครั้ง", "password not match", nil)
		return "", errors.New("password not match")
	}

	token, _ := utils.GenTokenMember(c, user)

	if err := s.userRepo.UpdateToken(user.ID, token); err != nil {
		fmt.Println("error get user", err)
		utils.Response(c, http.StatusInternalServerError, 500, "เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้ง", err.Error(), nil)
		return "", err
	}

	return token, nil
}

func (s *AuthenticationService) Register(c *gin.Context, payload domain.PayloadUser) error {
	user := domain.User{
		ID:       primitive.NewObjectID(),
		Username: payload.Username,
		Password: utils.MD5Hash(payload.Password),
		Active:   true,
		UpdateAt: time.Now(),
		CreateAt: time.Now(),
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		fmt.Println("error create user", err)
		utils.Response(c, http.StatusInternalServerError, 500, "เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้ง", err.Error(), nil)
		return err
	}

	return nil
}

func (s *AuthenticationService) ResetPassword(c *gin.Context, payload domain.PayloadResetPassword) error {
	user, err := s.userRepo.GetUserByName(payload.Username)
	if err != nil {
		fmt.Println("error get user", err)
		utils.Response(c, http.StatusInternalServerError, 500, "ไม่พบยูสเซอร์ กรุณาลองใหม่อีกครั้ง", err.Error(), nil)
		return err
	}

	if utils.MD5Hash(payload.OldPassword) != user.Password {
		fmt.Println("error password not match", err)
		utils.Response(c, http.StatusBadRequest, 400, "รหัสผ่านไม่ถูกต้อง กรุณาลองใหม่อีกครั้ง", "password not match", nil)
		return errors.New("password not match")
	}

	if payload.OldPassword == payload.NewPassword {
		fmt.Println("error old password and new password is match", err)
		utils.Response(c, http.StatusBadRequest, 400, "กรุณาลองใหม่อีกครั้ง", "error old password and new password is match", nil)
		return errors.New("error old password and new password is match")
	}

	if err := s.userRepo.UpdatePasswordUser(user.ID, payload.NewPassword); err != nil {
		fmt.Println("error update password user", err)
		utils.Response(c, http.StatusInternalServerError, 500, "", err.Error(), nil)
		return err
	}

	return nil
}


