package models

import "time"

type Book struct {
	ID        int
	Title     string
	Author    string
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
