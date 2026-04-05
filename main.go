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
	taskController := controllers.TaskController{DB: db}

	//router
	router := gin.Default()
	router.GET("/",func(ctx *gin.Context){
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!, welcome to Errandify API",
		})
	})

	//user router
	router.POST("/users/login", userController.Login)
	router.POST("/users", userController.CreateAccount)
	router.DELETE("/users/:id", userController.DeleteAccount)
	router.GET("/users/employee", userController.GetEmployees)

	//task router
	router.POST("/tasks", taskController.CreateTask)
	router.DELETE("/tasks/:id", taskController.DeleteTask)
	router.PATCH("/tasks/:id/submit", taskController.SubmitTask)
	router.PATCH("/tasks/:id/reject", taskController.RejectTask)

	//attachment router
	router.Static("/attachments", "./attachments")

	//server config
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":"+port)
}