package service

import "golang-ercommerce/models"


type AuthService interface {
	Login(user models.User) (string, error)
	Register(user models.User) (error)
}