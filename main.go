package main

import (
	"errandify/config"
	"errandify/models"
	"errandify/controllers"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main(){

	//database connection, auto migrate, create owner account
	db := config.DatabaseConnection()
	db.AutoMigrate(&models.User{}, &models.Task{})
	config.CreateOwnerAccount(db)

	//controller
	userController := controllers.UserController{DB: db}

	//router
	router := gin.Default()
	router.GET("/",func(ctx *gin.Context){
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!, welcome to Errandify API",
		})
	})

	router.POST("/users/login", userController.Login)
	router.POST("/users", userController.CreateAccount)
	router.DELETE("/users/:id", userController.DeleteAccount)

	router.Static("/attachment", "./attachment")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":"+port)
}