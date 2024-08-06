package data

type Book struct {
	ID       int    `json:"id"`
	UniqueID string `json:"unique_id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Genre    string `json:"genre"`
	Year     int    `json:"year"`
}
