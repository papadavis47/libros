package unit

import (
	"testing"
	"time"

	"github.com/papadavis47/libros/internal/models"
	"github.com/papadavis47/libros/internal/utils"
)

// TestFormatDate tests the date formatting utility function
// This function is used throughout the UI to display human-readable dates
// with ordinal suffixes (1st, 2nd, 3rd, etc.)
func TestFormatDate(t *testing.T) {
	tests := []struct {
		name     string
		date     time.Time
		expected string
	}{
		{
			name:     "January 1st",
			date:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: "January 1st, 2024",
		},
		{
			name:     "March 2nd",
			date:     time.Date(2023, 3, 2, 0, 0, 0, 0, time.UTC),
			expected: "March 2nd, 2023",
		},
		{
			name:     "May 3rd",
			date:     time.Date(2022, 5, 3, 0, 0, 0, 0, time.UTC),
			expected: "May 3rd, 2022",
		},
		{
			name:     "July 4th",
			date:     time.Date(2021, 7, 4, 0, 0, 0, 0, time.UTC),
			expected: "July 4th, 2021",
		},
		{
			name:     "December 11th (special case)",
			date:     time.Date(2020, 12, 11, 0, 0, 0, 0, time.UTC),
			expected: "December 11th, 2020",
		},
		{
			name:     "September 12th (special case)",
			date:     time.Date(2019, 9, 12, 0, 0, 0, 0, time.UTC),
			expected: "September 12th, 2019",
		},
		{
			name:     "April 13th (special case)",
			date:     time.Date(2018, 4, 13, 0, 0, 0, 0, time.UTC),
			expected: "April 13th, 2018",
		},
		{
			name:     "August 21st",
			date:     time.Date(2017, 8, 21, 0, 0, 0, 0, time.UTC),
			expected: "August 21st, 2017",
		},
		{
			name:     "October 22nd",
			date:     time.Date(2016, 10, 22, 0, 0, 0, 0, time.UTC),
			expected: "October 22nd, 2016",
		},
		{
			name:     "November 23rd",
			date:     time.Date(2015, 11, 23, 0, 0, 0, 0, time.UTC),
			expected: "November 23rd, 2015",
		},
		{
			name:     "February 28th",
			date:     time.Date(2014, 2, 28, 0, 0, 0, 0, time.UTC),
			expected: "February 28th, 2014",
		},
		{
			name:     "June 30th",
			date:     time.Date(2013, 6, 30, 0, 0, 0, 0, time.UTC),
			expected: "June 30th, 2013",
		},
		{
			name:     "January 31st",
			date:     time.Date(2012, 1, 31, 0, 0, 0, 0, time.UTC),
			expected: "January 31st, 2012",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.FormatDate(tt.date)
			if result != tt.expected {
				t.Errorf("FormatDate(%v) = %q, want %q", tt.date, result, tt.expected)
			}
		})
	}
}

// TestFormatBookType_WithEnum tests book type formatting with enum values
// This function handles the conversion from internal enum types to display strings
func TestFormatBookType_WithEnum(t *testing.T) {
	tests := []struct {
		name     string
		bookType models.BookType
		expected string
	}{
		{"paperback enum", models.Paperback, "Paperback"},
		{"hardback enum", models.Hardback, "Hardback"},
		{"audio enum", models.Audio, "Audio"},
		{"digital enum", models.Digital, "Digital"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.FormatBookType(tt.bookType)
			if result != tt.expected {
				t.Errorf("FormatBookType(%v) = %q, want %q", tt.bookType, result, tt.expected)
			}
		})
	}
}

// TestFormatBookType_WithString tests book type formatting with string values
// This function handles string inputs (case-insensitive) and formats them properly
func TestFormatBookType_WithString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"lowercase paperback", "paperback", "Paperback"},
		{"uppercase hardback", "HARDBACK", "Hardback"},
		{"mixed case audio", "AuDiO", "Audio"},
		{"lowercase digital", "digital", "Digital"},
		{"unknown type", "unknown", "Unknown"},
		{"empty string", "", ""},
		{"random string", "magazine", "Magazine"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.FormatBookType(tt.input)
			if result != tt.expected {
				t.Errorf("FormatBookType(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestFormatBookType_WithOtherTypes tests book type formatting with various input types
// This ensures the function gracefully handles unexpected input types
func TestFormatBookType_WithOtherTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"integer input", 123, "123"},
		{"boolean input", true, "true"},
		{"nil input", nil, "<nil>"},
		{"slice input", []string{"test"}, "[test]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.FormatBookType(tt.input)
			if result != tt.expected {
				t.Errorf("FormatBookType(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// BenchmarkFormatDate benchmarks the date formatting function
// This ensures the formatting performance is acceptable for UI rendering
func BenchmarkFormatDate(b *testing.B) {
	date := time.Date(2024, 6, 15, 12, 30, 45, 0, time.UTC)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utils.FormatDate(date)
	}
}

// BenchmarkFormatBookType benchmarks the book type formatting function
// This ensures the formatting performance is acceptable for list rendering
func BenchmarkFormatBookType(b *testing.B) {
	bookType := models.Paperback
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utils.FormatBookType(bookType)
	}
}