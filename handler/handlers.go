package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jerrylovee2/gogo/data"
)

func CreateBookHandler(w http.ResponseWriter, r *http.Request) {
	var newBook data.Book
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	data.InMemoryDB.Lock()
	defer data.InMemoryDB.Unlock()

	newBook.ID = data.InMemoryDB.NextBookID
	data.InMemoryDB.NextBookID++

	newBook.UniqueID = fmt.Sprintf("ID%d", newBook.ID)
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

func CreateMemberHandler(w http.ResponseWriter, r *http.Request) {
	var newMember data.Member
	err := json.NewDecoder(r.Body).Decode(&newMember)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	newMemberID := fmt.Sprintf("%03d", data.InMemoryDB.NextMemberID)
	newMember.ID = newMemberID

	data.InMemoryDB.Lock()
	defer data.InMemoryDB.Unlock()
	data.InMemoryDB.Members[newMember.ID] = newMember
	data.InMemoryDB.NextMemberID++

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newMember)
}

func GetMemberByIDHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	member, ok := data.InMemoryDB.Members[idParam]
	if !ok {
		http.Error(w, "Member not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(member)
}

func DeleteMemberByIDHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	_, ok := data.InMemoryDB.Members[idParam]
	if !ok {
		http.Error(w, "Member not found", http.StatusNotFound)
		return
	}

	data.InMemoryDB.Lock()
	defer data.InMemoryDB.Unlock()
	delete(data.InMemoryDB.Members, idParam)

	w.WriteHeader(http.StatusNoContent)
}

func CreateBorrowerHandler(w http.ResponseWriter, r *http.Request) {
	var newBorrower data.Borrower
	err := json.NewDecoder(r.Body).Decode(&newBorrower)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	data.InMemoryDB.Lock()
	defer data.InMemoryDB.Unlock()

	newBorrower.ID = data.InMemoryDB.NextBorrowerID
	data.InMemoryDB.Borrowers[newBorrower.ID] = newBorrower
	data.InMemoryDB.NextBorrowerID++

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newBorrower)
}

func GetBorrowerByIDHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	borrowerID, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid borrower ID", http.StatusBadRequest)
		return
	}

	data.InMemoryDB.RLock()
	defer data.InMemoryDB.RUnlock()

	borrower, ok := data.InMemoryDB.Borrowers[borrowerID]
	if !ok {
		http.Error(w, "Borrower not found", http.StatusNotFound)
		return
	}

	dueDate := time.Now().AddDate(0, 1, 0)
	penalty := 5.0

	borrowerInfo := struct {
		data.Borrower
		DueDate   time.Time `json:"due_date"`
		Penalty   float64   `json:"penalty_per_day"`
		Penalties float64   `json:"penalties"`
	}{
		Borrower:  borrower,
		DueDate:   dueDate,
		Penalty:   penalty,
		Penalties: calculatePenalties(dueDate, penalty),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(borrowerInfo)
}

func calculatePenalties(dueDate time.Time, penaltyPerDay float64) float64 {
	daysLate := time.Since(dueDate).Hours() / 24
	if daysLate <= 0 {
		return 0
	}
	return daysLate * penaltyPerDay
}

func DeleteBorrowerByIDHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	borrowerID, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid borrower ID", http.StatusBadRequest)
		return
	}

	data.InMemoryDB.Lock()
	defer data.InMemoryDB.Unlock()

	if _, ok := data.InMemoryDB.Borrowers[borrowerID]; !ok {
		http.Error(w, "Borrower not found", http.StatusNotFound)
		return
	}

	delete(data.InMemoryDB.Borrowers, borrowerID)

	w.WriteHeader(http.StatusNoContent)
}
