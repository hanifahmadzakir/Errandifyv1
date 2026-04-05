package main

import (
	"context"
	"errandify/config"
	"errandify/controllers"
	"errandify/models"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Database connection, auto migrate, create owner account
	db := config.DatabaseConnection()
	db.AutoMigrate(&models.User{}, &models.Task{})
	config.CreateOwnerAccount(db)

	// 2. Controller
	userController := controllers.UserController{DB: db}
	taskController := controllers.TaskController{DB: db}

	// 3. Router
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!, welcome to Errandify API",
		})
	})

	// User router
	router.POST("/users/login", userController.Login)
	router.POST("/users", userController.CreateAccount)
	router.DELETE("/users/:id", userController.DeleteAccount)
	router.GET("/users/employee", userController.GetEmployees)

	// Task router
	router.POST("/tasks", taskController.CreateTask)
	router.DELETE("/tasks/:id", taskController.DeleteTask)
	router.PATCH("/tasks/:id/submit", taskController.SubmitTask)
	router.PATCH("/tasks/:id/reject", taskController.RejectTask)

	// Attachment router
	router.Static("/attachments", "./attachments")

	// 4. Server config
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create HTTP server with custom configuration
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// 5. Run server in a goroutine so that it doesn't block the main thread and allows us to listen for shutdown signals
	go func() {
		log.Printf("Running on port %s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to run server: %s\n", err)
		}
	}()

	// 6. Graceful Shutdown Implementation
	// Create channel to listen for signals from the OS (like Ctrl+C or stop command from Docker)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Code will be blocked here until a signal is received in the 'quit' channel
	<-quit
	log.Println("Signal received for graceful shutdown...")

	// give the server 5 seconds to finish ongoing requests before forcefully shutting down
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shutdown the server. If it takes more than 5 seconds, it will forcefully shut down.
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server force shutdown: ", err)
	}

	// (close database connection
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
		log.Println("database connection closed.")
	}

	log.Println("Server Errandify shutdown successfully.")
}