package main

import (
	"context"
	"fmt"
	"github.com/upekZ/rest-api-go/internal/servermanager"
)

func main() {
	app := servermanager.New()

	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println("server failure: %w", err)
	}
}
