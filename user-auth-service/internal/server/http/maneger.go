package http

import (
	"awesomeProject4/user-auth-service/internal/domains/models"
	"awesomeProject4/user-auth-service/internal/usecase"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// UserHandler обрабатывает HTTP-запросы, связанные с пользователями
type UserHandler struct {
	UserUC *usecase.UserUseCase
}

// NewUserHandler создает новый экземпляр UserHandler
func NewUserHandler(userUC *usecase.UserUseCase) *UserHandler {
	return &UserHandler{UserUC: userUC}
}

// Register регистрирует нового пользователя
func (h *UserHandler) Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	if err := h.UserUC.RegisterUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	massage := fmt.Sprintf("Пользователь зарегистрирован успешно , ваш ID:%d", user.ID)
	c.JSON(http.StatusCreated, gin.H{"message": massage})
}

// ConfirmEmail подтверждает email пользователя
func (h *UserHandler) ConfirmEmail(c *gin.Context) {
	var input struct {
		UserID int    `json:"user_id" binding:"required"`
		Code   string `json:"code" binding:"required"`
	}
	logrus.Println("ConfirmEmail")
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Println("ShouldBindJSON")
	if err := h.UserUC.ConfirmEmail(c.Request.Context(), input.UserID, input.Code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logrus.Println("ConfirmEmail")
	c.JSON(http.StatusOK, gin.H{"message": "Email подтвержден успешно"})
}

// Login авторизует пользователя и возвращает JWT-токен
func (h *UserHandler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.UserUC.Login(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// ForgotPassword инициирует процесс восстановления пароля
func (h *UserHandler) ForgotPassword(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.UserUC.ForgotPassword(c.Request.Context(), input.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Инструкция по восстановлению пароля отправлена на email"})
}

// ResetPassword сбрасывает пароль пользователя
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var input struct {
		UserID      int    `json:"user_id" binding:"required"`
		Code        string `json:"code" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.UserUC.ResetPassword(c.Request.Context(), input.UserID, input.Code, input.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пароль сброшен успешно"})
}
