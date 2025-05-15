package main

import (
	"fmt"
	"github.com/upekZ/rest-api-go/internal/database"
	"github.com/upekZ/rest-api-go/internal/handlers"
	"github.com/upekZ/rest-api-go/internal/services"
)

func main() {
	dbConn, err := database.NewPostgresConn()

	if err != nil {
		fmt.Println("db connection failure: %w", err)
		return
	}

	userService := services.NewUserService(dbConn)

	app := handlers.NewServer(userService)

	if err = app.Start(); err != nil {

		fmt.Println("server failure: %w", err)
	}
}
