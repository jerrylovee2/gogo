package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jerrylovee2/gogo/data"
)

func CreateBookHandler(c *gin.Context) {
	var newBook data.Book
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, data.ErrorResponse{Error: "Invalid JSON"})
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

	c.JSON(http.StatusOK, newBook)
}

func DeleteBookHandler(c *gin.Context) {
	idParam := c.Query("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, data.ErrorResponse{Error: "Invalid book ID"})
		return
	}

	data.InMemoryDB.Lock()
	defer data.InMemoryDB.Unlock()
	if _, ok := data.InMemoryDB.Books[id]; !ok {
		c.JSON(http.StatusNotFound, data.ErrorResponse{Error: "Book not found"})
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

	c.Status(http.StatusNoContent)
}

func GetAllBooksHandler(c *gin.Context) {
	data.InMemoryDB.RLock()
	defer data.InMemoryDB.RUnlock()

	var books []data.Book
	for _, book := range data.InMemoryDB.Books {
		books = append(books, book)
	}

	c.JSON(http.StatusOK, books)
}

func SearchBooksHandler(c *gin.Context) {
	yearParam := c.Query("year")
	authorParam := c.Query("author")
	genreParam := c.Query("genre")

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

	c.JSON(http.StatusOK, filteredBooks)
}

func CreateMemberHandler(c *gin.Context) {
	var newMember data.Member
	if err := c.ShouldBindJSON(&newMember); err != nil {
		c.JSON(http.StatusBadRequest, data.ErrorResponse{Error: "Invalid JSON"})
		return
	}

	// Validate the length of the member's name
	if len(newMember.Name) > 15 {
		c.JSON(http.StatusBadRequest, data.ErrorResponse{Error: "Member name should not exceed 15 characters"})
		return
	}

	newMemberID := fmt.Sprintf("%03d", data.InMemoryDB.NextMemberID)
	newMember.ID = newMemberID

	data.InMemoryDB.Lock()
	defer data.InMemoryDB.Unlock()
	data.InMemoryDB.Members[newMember.ID] = newMember
	data.InMemoryDB.NextMemberID++

	c.JSON(http.StatusOK, newMember)
}

func GetMemberByIDHandler(c *gin.Context) {
	idParam := c.Query("id")
	member, ok := data.InMemoryDB.Members[idParam]
	if !ok {
		c.JSON(http.StatusNotFound, data.ErrorResponse{Error: "Member not found"})
		return
	}

	c.JSON(http.StatusOK, member)
}

func DeleteMemberByIDHandler(c *gin.Context) {
	idParam := c.Query("id")
	_, ok := data.InMemoryDB.Members[idParam]
	if !ok {
		c.JSON(http.StatusNotFound, data.ErrorResponse{Error: "Member not found"})
		return
	}

	data.InMemoryDB.Lock()
	defer data.InMemoryDB.Unlock()
	delete(data.InMemoryDB.Members, idParam)

	c.Status(http.StatusNoContent)
}

func CreateBorrowerHandler(c *gin.Context) {
	var newBorrower data.Borrower
	if err := c.ShouldBindJSON(&newBorrower); err != nil {
		c.JSON(http.StatusBadRequest, data.ErrorResponse{Error: "Invalid JSON"})
		return
	}

	data.InMemoryDB.Lock()
	defer data.InMemoryDB.Unlock()

	newBorrower.ID = data.InMemoryDB.NextBorrowerID
	data.InMemoryDB.Borrowers[newBorrower.ID] = newBorrower
	data.InMemoryDB.NextBorrowerID++

	c.JSON(http.StatusOK, newBorrower)
}

func GetBorrowerByIDHandler(c *gin.Context) {
	idParam := c.Query("id")
	borrowerID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, data.ErrorResponse{Error: "Invalid borrower ID"})
		return
	}

	data.InMemoryDB.RLock()
	defer data.InMemoryDB.RUnlock()

	borrower, ok := data.InMemoryDB.Borrowers[borrowerID]
	if !ok {
		c.JSON(http.StatusNotFound, data.ErrorResponse{Error: "Borrower not found"})
		return
	}

	dueDate := time.Now().AddDate(0, 1, 0)
	penalty := 5.0

	borrowerInfo := data.BorrowerInfo{
		Borrower:  borrower,
		DueDate:   dueDate,
		Penalty:   penalty,
		Penalties: calculatePenalties(dueDate, penalty),
	}

	c.JSON(http.StatusOK, borrowerInfo)
}

func calculatePenalties(dueDate time.Time, penaltyPerDay float64) float64 {
	daysLate := time.Since(dueDate).Hours() / 24
	if daysLate <= 0 {
		return 0
	}
	return daysLate * penaltyPerDay
}

func DeleteBorrowerByIDHandler(c *gin.Context) {
	idParam := c.Query("id")
	borrowerID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, data.ErrorResponse{Error: "Invalid borrower ID"})
		return
	}

	data.InMemoryDB.Lock()
	defer data.InMemoryDB.Unlock()

	if _, ok := data.InMemoryDB.Borrowers[borrowerID]; !ok {
		c.JSON(http.StatusNotFound, data.ErrorResponse{Error: "Borrower not found"})
		return
	}

	delete(data.InMemoryDB.Borrowers, borrowerID)

	c.Status(http.StatusNoContent)
}
