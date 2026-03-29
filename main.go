package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main(){
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