package models

import "time"

type BookType string

const (
	Paperback BookType = "paperback"
	Hardback  BookType = "hardback"
	Audio     BookType = "audio"
	Digital   BookType = "digital"
)

type Book struct {
	ID        int
	Title     string
	Author    string
	Type      BookType
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Screen int

const (
	MenuScreen Screen = iota
	AddBookScreen
	ListBooksScreen
	BookDetailScreen
	EditBookScreen
)
