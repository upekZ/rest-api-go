package main

import (
	"context"
	"fmt"

	"github.com/upekZ/rest-api-go/application"
)

func main() {
	app := application.New()

	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println("server failure: %w", err)
	}
}
