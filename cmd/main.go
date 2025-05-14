package main

import (
	"fmt"
	"github.com/upekZ/rest-api-go/internal/servermanager"
)

func main() {
	app := servermanager.New()

	err := app.Start()
	if err != nil {
		fmt.Println("server failure: %w", err)
	}
}
