package main

import (
	"fmt"
	"github.com/upekZ/rest-api-go/internal/cache"
	"github.com/upekZ/rest-api-go/internal/database"
	"github.com/upekZ/rest-api-go/internal/handlers"
	"github.com/upekZ/rest-api-go/internal/services"
	"github.com/upekZ/rest-api-go/internal/websocketService"
)

func main() {
	dbConn, err := database.NewPostgresConn()

	if err != nil {
		fmt.Printf("db connection failure: %v\n", err)
		return
	}

	userCache := cache.NewCache()

	hub := websocketService.NewHub()

	if hub == nil {
		fmt.Printf("web socket initialization failed\n")
		return
	}

	go hub.Run()

	userService := services.NewUserService(dbConn, userCache, hub)

	app := handlers.NewServer(userService)

	fmt.Printf("web server initialization success\n")

	if err = app.Start(); err != nil {
		fmt.Printf("server failure: %v\n", err)
	}
}
