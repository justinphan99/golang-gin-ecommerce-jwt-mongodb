package service

import (
	"golang-ercommerce/helpers"
	"golang-ercommerce/models"
	"golang-ercommerce/repositories"

	"github.com/go-playground/validator/v10"
)

type AuthServiceImpl struct {
	UserRepository repositories.UserRepository
	Validate       *validator.Validate
}

func NewAuthServiceImpl(userRepository repositories.UserRepository, validate *validator.Validate) AuthService {
	return &AuthServiceImpl{
		UserRepository: userRepository,
		Validate:       validate,
	}
}


// Register implements AuthService.
func (a *AuthServiceImpl) Register(user models.User) error {
	err := a.UserRepository.Register(user)
	return err
}

func (a *AuthServiceImpl) Login(user models.User) (string, error) {
	founduser, err := a.UserRepository.Login(user)
	if err != nil {
		return "", err
	}
	PasswordIsValid, _ := helpers.VerifyPassword(*user.Password, *founduser.Password)

	if !PasswordIsValid {
		return "", err
	}
	token, refreshToken, _ := helpers.TokenGenerator(*founduser.Email, *founduser.First_Name, *founduser.Last_Name, founduser.User_ID)
	helpers.UpdateAllTokens(token, refreshToken, founduser.User_ID)
	return *founduser.Token, nil
}
