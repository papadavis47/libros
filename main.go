package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/papadavis47/libros/internal/database"
	"github.com/papadavis47/libros/internal/ui"
)

func main() {
	db, err := database.New("books.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	model := ui.NewModel(db)
	p := tea.NewProgram(model)
	
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
