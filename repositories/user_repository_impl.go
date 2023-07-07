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

type UserRepositoryImpl struct {
	UserCollection *mongo.Collection
}

func NewUsersRepositoryImpl(userCollection *mongo.Collection) UserRepository {
	return &UserRepositoryImpl{
		UserCollection: userCollection,
	}
}

// AddProductToUserCart implements UserRepository.
func (u *UserRepositoryImpl) AddProductToUserCart(productcart []models.ProductUser, userID string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return helpers.ErrUserIDIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productcart}}}}}}
	_, err = u.UserCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return helpers.ErrCantUpdateUser
	}
	return nil
}

// RemoveCartItemByProductIDAndUserID implements UserRepository.
func (u *UserRepositoryImpl) RemoveCartItemByProductIDAndUserID(productID primitive.ObjectID, userID string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return helpers.ErrUserIDIsNotValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productID}}}
	_, err = u.UserCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return helpers.ErrCantRemoveItem
	}
	return nil
}

// GetUserByUserID implements UserRepository. 
func (u *UserRepositoryImpl) GetUserByUserID(userID primitive.ObjectID) (models.User, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var filleduser models.User
	err := u.UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: userID}}).Decode(&filleduser)
	if err != nil {
		log.Println(err)
		return filleduser, helpers.ErrUserNotFound
	}

	return filleduser, nil
}

// GetUserCartByUserID implements UserRepository. 
func (u *UserRepositoryImpl) GetUserCartByUserID(userID primitive.ObjectID) ([]bson.M, models.User, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var filleduser models.User
	var listing []bson.M

	filleduser, err := u.GetUserByUserID(userID)
	if err != nil {
		log.Println(err)
		return listing, filleduser, err
	}
	
	filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: userID}}}}
	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}
	pointcursor, err := u.UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})
	if err != nil {
		log.Println(err)
	}
	if err = pointcursor.All(ctx, &listing); err != nil {
		log.Println(err)
		return listing, filleduser, err
	}
	return listing, filleduser, nil
}

// Register implements UserRepository. 
func (u *UserRepositoryImpl) Register(user models.User) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	count, err := u.UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		log.Panic(err)
		return err
	}
	if count > 0 {
		return helpers.ErrUserAlreadyExists
	}

	count, err = u.UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
	if err != nil {
		log.Panic(err)
		return err
	}
	if count > 0 {
		return helpers.ErrPhoneIsAlreadyInUse
	}

	password, err := helpers.HashPassword(*user.Password)
	if err != nil {
		log.Panic(err)
		return err
	}
	user.Password = &password
	user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.User_ID = user.ID.Hex()
	token, refreshtoken, _ := helpers.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)
	user.Token = &token
	user.Refresh_Token = &refreshtoken
	user.UserCart = make([]models.ProductUser, 0)
	user.Address_Details = make([]models.Address, 0)
	user.Order_Status = make([]models.Order, 0)
	_, inserterr := u.UserCollection.InsertOne(ctx, user)
	return inserterr
}

func (u *UserRepositoryImpl) Login(user models.User) (models.User, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var founduser models.User
	err := u.UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
	if err != nil {
		log.Print(err)
		return founduser, err
	}
	return founduser, nil
}