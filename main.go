package golangecommerce

import (
	"golang-ercommerce/config"
	"golang-ercommerce/controllers"
	"golang-ercommerce/database"
	"golang-ercommerce/helpers"
	"golang-ercommerce/repositories"
	"golang-ercommerce/routes"
	"golang-ercommerce/service"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {

	// Load config
	loadConfig, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load enviroment variables", err)
	}

	serverPort := loadConfig.ServerPort
	if serverPort == "" {
		serverPort = ":8000"
	}
	// Database
	// db := config.ConnectionDB(&loadConfig)
	var Client *mongo.Client = database.DBSet()
	UsersCollection := database.UserCollection(Client)
	ProductCollection := database.ProductCollection(Client)

	// db.Table("users").AutoMigrate(&model.Users{})

	// Init repository
	// usersRepository := repository.NewUsersRepositoryImpl(db)
	// productRepository := repository.NewProductRepositoryImpl(db)
	usersRepository := repositories.NewUsersRepositoryImpl(UsersCollection)
	productRepository := repositories.NewProductRepositoryImpl(ProductCollection)

	// Init service
	// authenticationService := service.NewAuthenticationServiceImpl(usersRepository, validate)
	// productService := service.NewProductServiceImpl(productRepository, validate)
	validate := validator.New()
	cartService := service.NewCartServiceImpl(productRepository, usersRepository, validate)

	// Init controller
	// authenticationController := controller.NewAuthenticationController(authenticationService)
	// usersController := controller.NewUsersController(usersRepository)
	// productController := controller.NewProductController(productService)
	cartController := controllers.NewCartController(cartService)

	routes := routes.NewRoutes(cartController)

	server := &http.Server{
		Addr:    serverPort,
		Handler: routes,
	}

	server_err := server.ListenAndServe()
	helpers.ErrorPanic(server_err)
}
