package messages

import "github.com/papadavis47/libros/internal/models"

type SaveMsg struct {
	Err error
}

type UpdateMsg struct {
	Err error
}

type DeleteMsg struct {
	Err error
}

type LoadBooksMsg struct {
	Books []models.Book
	Err   error
}
