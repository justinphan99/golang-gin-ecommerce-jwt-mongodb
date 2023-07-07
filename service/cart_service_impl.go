package service

import (
	"golang-ercommerce/models"
	"golang-ercommerce/repositories"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartServiceImpl struct {
	ProductRepository repositories.ProductRepository
	UserRepository    repositories.UserRepository
	Validate          *validator.Validate
}

func NewCartServiceImpl(productRepository repositories.ProductRepository,
	userRepository repositories.UserRepository,
	validate *validator.Validate) CartService {
	return &CartServiceImpl{
		ProductRepository: productRepository,
		UserRepository:    userRepository,
		Validate:          validate,
	}
}

// AddProductToCart implements CartService.
func (c *CartServiceImpl) AddProductToCart(productID primitive.ObjectID, userID string) error {
	productcart, err := c.ProductRepository.Find(productID)
	if err != nil {
		return err
	}
	err = c.UserRepository.AddProductToUserCart(productcart, userID)
	return err
}

// RemoveCartItem implements CartService.
func (c *CartServiceImpl) RemoveCartItem(productID primitive.ObjectID, userID string) error {
	err := c.UserRepository.RemoveCartItemByProductIDAndUserID(productID, userID)
	return err 
}

func (c *CartServiceImpl) GetItemFromCart(userID primitive.ObjectID) ([]bson.M, models.User, error) {
	listing, filleduser, err := c.UserRepository.GetUserCartByUserID(userID)
	if err != nil {
		return listing, filleduser, err
	}
	
	return listing, filleduser, nil
}
