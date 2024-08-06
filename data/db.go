package data

import (
	"sync"
)

// InMemoryDB simulates an in-memory database
var InMemoryDB = struct {
	Books          map[int]Book
	Members        map[string]Member
	Borrowers      map[int]Borrower
	Indices        map[string]map[string][]int
	NextBookID     int
	NextMemberID   int
	NextBorrowerID int
	sync.RWMutex
}{Books: make(map[int]Book), Members: make(map[string]Member), Borrowers: make(map[int]Borrower), Indices: make(map[string]map[string][]int)}
