package main

import (
	"errandify/config"
	"errandify/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main(){

	//database connection, auto migrate, create owner account
	db := config.DatabaseConnection()
	db.AutoMigrate(&models.User{}, &models.Task{})
	config.CreateOwnerAccount(db)

	//router
	router := gin.Default()
	router.GET("/",func(ctx *gin.Context){
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!, welcome to Errandify API",
		})
	})
	router.Static("/attachment", "./attachment")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":"+port)
}