package models

import (
	"log"

	"github.com/subodh/library/pkg/config"
)

func GetUserStats() *[]Book {
	// Query the data
	DB := config.GetDB()

	// Use DB.Raw() to execute a raw SQL query
	rows, err := DB.Raw("SELECT Id, name,author, publication,created_at,updated_at,deleted_at,pdf_data FROM bookstore.books ORDER BY Id DESC").Rows()
	if err != nil {
		log.Fatalf("Error executing raw query: %v", err)
	}
	defer rows.Close()

	// Slice to hold the user stats
	var books []Book

	// Iterate over the rows
	for rows.Next() {
		var book Book
		// Scan the row into the struct fields
		err := rows.Scan(
			&book.ID,          // ID from gorm.Model
			&book.Name,        // Name from Book
			&book.Author,      // Author from Book
			&book.Publication, // Publication from Book
			&book.CreatedAt,
			&book.UpdatedAt,
			&book.DeletedAt,
			&book.PDFData,
		)

		if err != nil {
			log.Fatal(err)
		}

		// Append to the slice
		books = append(books, book)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return &books
}
