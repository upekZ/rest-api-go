package main

import (
	"fmt"
	"github.com/upekZ/rest-api-go/internal/cache"
	"github.com/upekZ/rest-api-go/internal/database"
	"github.com/upekZ/rest-api-go/internal/handlers"
	"github.com/upekZ/rest-api-go/internal/services"
)

func main() {
	dbConn, err := database.NewPostgresConn()

	if err != nil {
		fmt.Printf("db connection failure: %v\n", err)
		return
	}

	userCache := cache.NewCache()
	userService := services.NewUserService(dbConn, userCache)

	app := handlers.NewServer(userService)

	if err = app.Start(); err != nil {

		fmt.Println("server failure: %v", err)
	}
}
