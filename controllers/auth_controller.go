package controllers

import (
	"golang-ercommerce/models"
	"golang-ercommerce/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

func (controller *AuthController) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := controller.AuthService.Register(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		c.JSON(http.StatusCreated, "Successfully Signed Up!!")
	}
}

func (controller *AuthController) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		token, err := controller.AuthService.Login(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		c.JSON(http.StatusFound, gin.H{"token": token})
	}
}