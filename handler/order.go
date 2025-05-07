package handler

import (
	"fmt"
	"net/http"
)

type Order struct{}

func (o *Order) Create(writer http.ResponseWriter, reader *http.Request) {
	fmt.Println("Order creation requested")
}

func (o *Order) List(writer http.ResponseWriter, reader *http.Request) {
	fmt.Println("Order Listing")
}

func (o *Order) GetByID(writer http.ResponseWriter, reader *http.Request) {
	fmt.Println("Order by ID")
}

func (o *Order) UpdateByID(writer http.ResponseWriter, reader *http.Request) {
	fmt.Println("Order update by ID")
}

func (o *Order) DeleteByID(writer http.ResponseWriter, reader *http.Request) {
	fmt.Println("Order Delete by ID")
}
