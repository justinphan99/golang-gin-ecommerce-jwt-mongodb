package repositories

import (
	"context"
	"golang-ercommerce/helpers"
	"golang-ercommerce/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepositoryImpl struct {
	Collection *mongo.Collection
}

func NewProductRepositoryImpl(Collection *mongo.Collection) ProductRepository {
	return &ProductRepositoryImpl{Collection: Collection}
}

// Find implements ProductRepository.
func (p *ProductRepositoryImpl) Find(productID primitive.ObjectID) ([]models.ProductUser, error) {
	var productcart []models.ProductUser
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	searchfromdb, err := p.Collection.Find(ctx, bson.M{"_id": productID})
	if err != nil {
		log.Println(err)
		return productcart, helpers.ErrCantFindProduct
	}
	err = searchfromdb.All(ctx, &productcart)
	if err != nil {
		log.Println(err)
		return productcart, helpers.ErrCantDecodeProducts
	}
	return productcart, nil
}


