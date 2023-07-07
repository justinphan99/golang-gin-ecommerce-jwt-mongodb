package service

import (
	"golang-ercommerce/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartService interface {
	AddProductToCart(productID primitive.ObjectID, userID string) error
	RemoveCartItem(productID primitive.ObjectID, userID string) error
	GetItemFromCart(userID primitive.ObjectID) ([]bson.M, models.User, error)
	// RemoveCartItem(prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error 	
	// BuyItemFromCart(userCollection *mongo.Collection, userID string) error 
	// InstantBuyer(prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, UserID string) error
}