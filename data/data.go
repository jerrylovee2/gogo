package data

import (
	"sync"
)

type Member struct {
	ID          int    `json:"id"`
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
	MemberID int   `json:"member_id"`
	Books    []int `json:"books"`
}

var InMemoryDB struct {
	sync.RWMutex
	Books     map[int]Book
	NextID    int
	Indices   map[string]map[string][]int
	Members   map[int]Member
	Borrowers map[int]Borrower
}

func Initialize() {
	InMemoryDB.Books = make(map[int]Book)
	InMemoryDB.Indices = make(map[string]map[string][]int)
	InMemoryDB.NextID = 1
}
