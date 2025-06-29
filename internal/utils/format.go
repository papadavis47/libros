package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/papadavis47/libros/internal/models"
)

// FormatDate formats a time.Time value into a human-readable date string
// with the month name, day with ordinal suffix, and year.
//
// Parameters:
//   - t: The time.Time value to format
//
// Returns:
//   - string: Formatted date with month name, day with ordinal suffix, and year
func FormatDate(t time.Time) string {
	day := t.Day()
	var suffix string
	// Determine appropriate ordinal suffix for the day
	switch {
	case day >= 11 && day <= 13:
		// Special case: 11th, 12th, 13th (not 11st, 12nd, 13rd)
		suffix = "th"
	case day%10 == 1:
		suffix = "st" // 1st, 21st, 31st
	case day%10 == 2:
		suffix = "nd" // 2nd, 22nd
	case day%10 == 3:
		suffix = "rd" // 3rd, 23rd
	default:
		suffix = "th" // 4th, 5th, 6th, 7th, 8th, 9th, 10th, etc.
	}
	return fmt.Sprintf("%s %d%s, %d", t.Format("January"), day, suffix, t.Year())
}

// FormatBookType converts BookType enum values to capitalized display names
// Handles both models.BookType enum and string representations
func FormatBookType(bookType interface{}) string {
	switch v := bookType.(type) {
	case models.BookType:
		switch v {
		case models.Paperback:
			return "Paperback"
		case models.Hardback:
			return "Hardback"
		case models.Audio:
			return "Audio"
		case models.Digital:
			return "Digital"
		default:
			return string(v)
		}
	case string:
		switch strings.ToLower(v) {
		case "paperback":
			return "Paperback"
		case "hardback":
			return "Hardback"
		case "audio":
			return "Audio"
		case "digital":
			return "Digital"
		default:
			// Fallback: capitalize first letter for unknown types
			if len(v) == 0 {
				return v
			}
			return strings.ToUpper(string(v[0])) + strings.ToLower(v[1:])
		}
	default:
		return fmt.Sprintf("%v", bookType)
	}
}