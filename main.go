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

	http.HandleFunc("/members/create", handlers.CreateMemberHandler)
	http.HandleFunc("/members/get", handlers.GetMemberByIDHandler)
	http.HandleFunc("/members/delete", handlers.DeleteMemberByIDHandler)

	http.HandleFunc("/borrowers/create", handlers.CreateBorrowerHandler)
	http.HandleFunc("/borrowers/get", handlers.GetBorrowerByIDHandler)
	http.HandleFunc("/borrowers/delete", handlers.DeleteBorrowerByIDHandler)

	port := ":8081"
	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
