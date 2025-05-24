package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"apiserver/internal/generated/api" // Generated API server
	"apiserver/internal/handlers"
	"apiserver/internal/repositories"
	"apiserver/internal/usecases"
	// "apiserver/internal/routes" // Old routes, replaced by api.RegisterHandlers
)

func main() {
	// Database Connection
	mysqlUser := os.Getenv("MYSQL_USER")
	if mysqlUser == "" {
		mysqlUser = "user" // Default
		log.Println("Warning: MYSQL_USER not set, using default 'user'")
	}
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	if mysqlPassword == "" {
		mysqlPassword = "password" // Default
		log.Println("Warning: MYSQL_PASSWORD not set, using default 'password'")
	}
	mysqlHost := os.Getenv("MYSQL_HOST")
	if mysqlHost == "" {
		mysqlHost = "127.0.0.1" // Default
		log.Println("Warning: MYSQL_HOST not set, using default '127.0.0.1'")
	}
	mysqlPort := os.Getenv("MYSQL_PORT")
	if mysqlPort == "" {
		mysqlPort = "3306" // Default
		log.Println("Warning: MYSQL_PORT not set, using default '3306'")
	}
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")
	if mysqlDatabase == "" {
		mysqlDatabase = "apidb" // Default
		log.Println("Warning: MYSQL_DATABASE not set, using default 'apidb'")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase)

	dbConn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer dbConn.Close()

	err = dbConn.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to the database.")

	// Initialize layers
	userRepo := repositories.NewUserRepository(dbConn)
	userInteractor := usecases.NewUserInteractor(userRepo)
	// UserHandler implements api.ServerInterface
	userHandler := handlers.NewUserHandler(userInteractor) 

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register handlers - oapi-codegen generates this function
	// The first argument is the Echo instance, the second is our ServerInterface implementation
	api.RegisterHandlers(e, userHandler) 

	// The old routes.SetupRoutes(e) is now replaced by api.RegisterHandlers(e, userHandler)

	// Start server
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080" // Default port
		log.Println("Warning: SERVER_PORT not set, using default '8080'")
	}
	log.Printf("Starting server on :%s", serverPort)
	if err := e.Start(":" + serverPort); err != nil {
		e.Logger.Fatal(err)
	}
}
