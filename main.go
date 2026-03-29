package main

import (
	"errandify/config"
	"errandify/models"
	"net/http"

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
	router.Run(":8080")
}