package routes

import (
	"apiserver/internal/handlers" // Assuming 'apiserver' is the module name

	"github.com/labstack/echo/v4"
)

// SetupRoutes defines all the application routes.
func SetupRoutes(e *echo.Echo) {
	// User routes based on OpenAPI specification
	// GET /v1/users -> maps to getUsers operationId
	e.GET("/v1/users", handlers.GetUsers)

	// POST /v1/user -> maps to post-user operationId
	e.POST("/v1/user", handlers.CreateUser)

	// PATCH /v1/users/{user_id} -> maps to path-user operationId
	e.PATCH("/v1/users/:user_id", handlers.UpdateUser) // Echo uses :param for path parameters

	// DELETE /v1/users/{user_id} -> maps to delete-user operationId
	e.DELETE("/v1/users/:user_id", handlers.DeleteUser)
}
