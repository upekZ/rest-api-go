package main

import (
	"fmt"
	"github.com/upekZ/rest-api-go/api/handler"
	"github.com/upekZ/rest-api-go/database"
	"github.com/upekZ/rest-api-go/internal/servermanager"
)

func main() {
	dbConn, err := database.NewPostgresConn()

	reqHandler := handler.NewHandler(dbConn)

	app := servermanager.NewServer(reqHandler)

	err = app.Start()
	if err != nil {
		fmt.Println("server failure: %w", err)
	}
}
