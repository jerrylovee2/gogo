package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jerrylovee2/gogo/data"
	handlers "github.com/jerrylovee2/gogo/handler"
)

func main() {
	data.Initialize()

	http.HandleFunc("/books/create", handlers.CreateBookHandler)
	http.HandleFunc("/books/delete", handlers.DeleteBookHandler)
	http.HandleFunc("/books/all", handlers.GetAllBooksHandler)
	http.HandleFunc("/books/search", handlers.SearchBooksHandler)

	port := ":8081"
	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
