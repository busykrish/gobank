package routes

import (
	"net/http"

	controllers "github.com/cavdy-play/go_db/controllers"
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	router.GET("/", welcome)
	router.GET("/customer", controllers.GetAllAccounts)
	router.POST("/createAccount", controllers.CreateAccount)
	router.POST("/createCustomer", controllers.CreateCustomer) // For creating customer
	router.GET("/getsingle", controllers.GetSingleAccount)
	router.PUT("/update", controllers.EditAccount)
	router.DELETE("/delete", controllers.DeleteAccount)
	router.NoRoute(notFound)
}

func welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Welcome To Krishna's API",
	})
	return
}

func notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status":  404,
		"message": "Route Not Found",
	})
	return
}
