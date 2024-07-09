package data

import (
	"sync"
)

type Member struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type Book struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Year     int    `json:"year"`
	Genre    string `json:"genre"`
	Edition  string `json:"edition"`
	UniqueID string `json:"unique_id"`
}

type Borrower struct {
	ID       int    `json:"id"`
	MemberID string `json:"member_id"`
	Books    []int  `json:"books"`
}

var InMemoryDB struct {
	sync.RWMutex
	Books          map[int]Book
	NextBookID     int
	Indices        map[string]map[string][]int
	Members        map[string]Member
	NextMemberID   int
	Borrowers      map[int]Borrower
	NextBorrowerID int
}

func Initialize() {
	InMemoryDB.Books = make(map[int]Book)
	InMemoryDB.Indices = make(map[string]map[string][]int)
	InMemoryDB.Members = make(map[string]Member)
	InMemoryDB.Borrowers = make(map[int]Borrower)
	InMemoryDB.NextBookID = 1
	InMemoryDB.NextMemberID = 1
	InMemoryDB.NextBorrowerID = 1
}
