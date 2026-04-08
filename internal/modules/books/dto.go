package books

import "time"

type CreateBookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

type UpdateBookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

type BookResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Year      int       `json:"year"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FindAllParams struct {
	Author string
	Title  string
	Page   int
	Limit  int
}
