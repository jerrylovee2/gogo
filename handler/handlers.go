package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jerrylovee2/gogo/data"
)

func CreateBookHandler(w http.ResponseWriter, r *http.Request) {
	var newBook data.Book
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	newBook.ID = data.InMemoryDB.NextID
	data.InMemoryDB.NextID++

	newBook.UniqueID = fmt.Sprintf("ID%d", newBook.ID)

	data.InMemoryDB.Lock()
	defer data.InMemoryDB.Unlock()
	data.InMemoryDB.Books[newBook.ID] = newBook

	if data.InMemoryDB.Indices[newBook.Genre] == nil {
		data.InMemoryDB.Indices[newBook.Genre] = make(map[string][]int)
	}
	data.InMemoryDB.Indices[newBook.Genre][newBook.Author] = append(data.InMemoryDB.Indices[newBook.Genre][newBook.Author], newBook.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newBook)
}

func DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	data.InMemoryDB.Lock()
	defer data.InMemoryDB.Unlock()
	if _, ok := data.InMemoryDB.Books[id]; !ok {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	delete(data.InMemoryDB.Books, id)

	for genre := range data.InMemoryDB.Indices {
		for author, ids := range data.InMemoryDB.Indices[genre] {
			var updatedIDs []int
			for _, bookID := range ids {
				if bookID != id {
					updatedIDs = append(updatedIDs, bookID)
				}
			}
			data.InMemoryDB.Indices[genre][author] = updatedIDs
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetAllBooksHandler(w http.ResponseWriter, r *http.Request) {
	data.InMemoryDB.RLock()
	defer data.InMemoryDB.RUnlock()

	var books []data.Book
	for _, book := range data.InMemoryDB.Books {
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func SearchBooksHandler(w http.ResponseWriter, r *http.Request) {
	yearParam := r.URL.Query().Get("year")
	authorParam := r.URL.Query().Get("author")
	genreParam := r.URL.Query().Get("genre")

	data.InMemoryDB.RLock()
	defer data.InMemoryDB.RUnlock()

	var filteredBooks []data.Book
	for _, book := range data.InMemoryDB.Books {
		if (yearParam == "" || strconv.Itoa(book.Year) == yearParam) &&
			(authorParam == "" || strings.Contains(strings.ToLower(book.Author), strings.ToLower(authorParam))) &&
			(genreParam == "" || strings.Contains(strings.ToLower(book.Genre), strings.ToLower(genreParam))) {
			filteredBooks = append(filteredBooks, book)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredBooks)
}
