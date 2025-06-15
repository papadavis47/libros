package models

type Book struct {
	ID     int
	Title  string
	Author string
}

type Screen int

const (
	MenuScreen Screen = iota
	AddBookScreen
	ListBooksScreen
	BookDetailScreen
	EditBookScreen
)
