package data

import "time"

type Borrower struct {
	ID       int       `json:"id"`
	MemberID string    `json:"member_id"`
	BookID   int       `json:"book_id"`
	Borrowed time.Time `json:"borrowed"`
}

type BorrowerInfo struct {
	Borrower
	DueDate   time.Time `json:"due_date"`
	Penalty   float64   `json:"penalty_per_day"`
	Penalties float64   `json:"penalties"`
}
