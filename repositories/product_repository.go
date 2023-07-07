package repositories

import (
	"golang-ercommerce/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductRepository interface {
	Find(productID primitive.ObjectID) ([]models.ProductUser, error)
}