package routes

import (
	"golang-ercommerce/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRoutes(cartController *controllers.CartController) *gin.Engine {

	service := gin.Default()

	service.GET("", func(context *gin.Context) {
		context.JSON(http.StatusOK, "welcome home")
	})

	service.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	router := service.Group("/api")
	{
		authenticationRouter := router.Group("/auth")
		{
			authenticationRouter.POST("/register")
			authenticationRouter.POST("/login")
		}
		// usersRouter := router.Group("/users")
		// {
		// 	usersRouter.GET("")
		// }
		// productRouter := router.Group("/product")
		// {
		// 	productRouter.GET("", controllers.SearchProduct())
		// 	productRouter.GET("/search", controllers.SearchProductByQuery())
		// 	productRouter.POST("/add", controllers.AddNewProduct())
		// }
		cartRouter := router.Group("/cart")
		{
			cartRouter.POST("/addtocart", cartController.AddToCart())
			cartRouter.POST("/removeitem", cartController.RemoveItem())
			cartRouter.GET("/listcart", cartController.GetItemFromCart())
		}
	}

	return service
}
