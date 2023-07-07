package repositories

import (
	"golang-ercommerce/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	AddProductToUserCart(productcart []models.ProductUser, userID string) error
	RemoveCartItemByProductIDAndUserID(productID primitive.ObjectID, userID string) error
	GetUserByUserID(userID primitive.ObjectID) (models.User, error)
	GetUserCartByUserID(userID primitive.ObjectID) ([]bson.M, models.User, error)
	Register(user models.User) error
	Login(user models.User) (models.User, error)
}