package controller

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/subodh/library/pkg/config"
	models "github.com/subodh/library/pkg/models"
)

// Define the data structure to pass to the template
type PageData struct {
	Title     string
	BookStats []models.Book
}

// handler function to respond to GET requests
func ReportHandler(w http.ResponseWriter, r *http.Request) {

	// Define the data to pass to the template
	data := PageData{
		Title:     "BOOKSTORE DATA",
		BookStats: models.GetAllBook(),
	}
	// Define the path to the template
	templatePath := "D:/practice/golang/devlelopment/newproject/library/views/home.html"
	// Ensure the template exists
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		log.Fatalf("Template file does not exist: %v", err)
	}
	// Parse and execute the template
	tmpl := template.Must(template.ParseFiles(templatePath))
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error rendering template for %v: %v", data.Title, err)
		return
	}
}

func ReportCSVHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write the header
	writer.Write([]string{"Id", "Name", "Author", "Publication", "Created_at", "Updated_at", "Deleted_at"})

	BookStats := models.GetUserStats()
	for _, user := range *BookStats {
		Deleted_at := "null"               // default to "null" if DeletedAt is nil
		if !user.DeletedAt.Time.IsZero() { // Check if DeletedAt is not zero
			// Format DeletedAt to a string (e.g., "2006-01-02 15:04:05")
			Deleted_at = user.DeletedAt.Time.Format("2006-01-02 15:04:05")
		}
		// Writing the user data to the CSV or output
		writer.Write([]string{
			strconv.Itoa(int(user.ID)), // Accessing the embedded ID field
			user.Name,
			user.Author,
			user.Publication,
			user.CreatedAt.Format("2006-01-02 15:04:05"),
			user.UpdatedAt.Format("2006-01-02 15:04:05"),
			//	user.DeletedAt.Time.Format("2006-01-02 15:04:05"),
			Deleted_at,
		})
	}

	writer.Flush()

	// Set the response headers
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=book_record.csv")
	if _, err := w.Write(buf.Bytes()); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
func UploadCSVHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form to retrieve file data
	err := r.ParseMultipartForm(0 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to parse form: %v", err), http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, _, err := r.FormFile("csvfile")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving the file: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()
	// Optional: Read the CSV content
	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Insert into DB using GORM
	for {
		records, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break // End of file reached
			}
		}

		// Assuming CSV format: Name, Author, Publication, CreatedAt, UpdatedAt
		book := models.Book{
			Name:        records[0],
			Author:      records[1],
			Publication: records[2],
		}

		// Insert the record into the database
		db := config.GetDB()
		if err := db.Create(&book).Error; err != nil {
			return
		}
	}
}

// Define the data structure to pass to the template
type PageSearchData struct {
	Title     string
	BookStats []models.Book
}

func SearchBookByName(w http.ResponseWriter, r *http.Request) {

	db := config.GetDB()
	name := r.URL.Query().Get("name")
	// Call the model function to search for books by name
	books, err := models.SearchBooksByName(db, name)
	if err != nil {
		log.Printf("Error searching books: %v", err)
		http.Error(w, "Error searching books", http.StatusInternalServerError)
		return
	}

	data := PageSearchData{
		Title:     "BookSearchData",
		BookStats: books,
	}
	// Define the path to the template
	templatePath := "D:/practice/golang/devlelopment/newproject/library/views/home.html"
	// Ensure the template exists
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		log.Fatalf("Template file does not exist: %v", err)
	}
	// Parse and execute the template
	tmpl := template.Must(template.ParseFiles(templatePath))
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error rendering template for %v: %v", data, err)
		return
	}
}

// UploadPDFHandler handles the upload of a PDF file for a specific book.
func UploadPDFHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form to retrieve file data
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to parse form: %v", err), http.StatusBadRequest)
		return
	}

	// Get the book ID from the form
	bookID := r.FormValue("book_id")
	if bookID == "" {
		http.Error(w, "Book ID is required", http.StatusBadRequest)
		return
	}
	// Get the PDF file from the form
	file, _, err := r.FormFile("pdffile")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving the file: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the PDF file into a byte slice
	pdfData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading PDF file: %v", err), http.StatusInternalServerError)
		return
	}

	// Update the book record in the database with the PDF data
	var book models.Book
	db := config.GetDB()
	if err := db.Where("id=?", bookID).First(&book).Error; err != nil {
		http.Error(w, fmt.Sprintf("Book not found: %v", err), http.StatusNotFound)
		return
	}

	book.PDFData = pdfData // Update the PDF data
	if err := db.Save(&book).Error; err != nil {
		http.Error(w, fmt.Sprintf("Error saving book record: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "PDF successfully uploaded for book ID %s", bookID)
}

func DownloadPDFHandler(w http.ResponseWriter, r *http.Request) {
	// Get the book ID from the query parameters
	bookIDStr := strings.TrimPrefix(r.URL.Path, "/download/")
	if bookIDStr == "" {
		http.Error(w, "Book ID is required", http.StatusBadRequest)
		return
	}

	// Convert book ID to uint
	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid Book ID", http.StatusBadRequest)
		return
	}

	// Fetch the book record from the database
	var book models.Book
	db := config.GetDB()
	if err := db.Where("id=?", bookID).First(&book).Error; err != nil {
		http.Error(w, fmt.Sprintf("Book not found: %v", err), http.StatusNotFound)
		return
	}
	// Check if PDF data exists
	if len(book.PDFData) == 0 {
		http.Error(w, "No PDF data available for this book", http.StatusNotFound)
		return
	}

	// Set headers to download the PDF
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.pdf\"", book.Name))
	w.Header().Set("Content-Length", strconv.Itoa(len(book.PDFData)))

	// Write the PDF data to the response
	if _, err := w.Write(book.PDFData); err != nil {
		log.Printf("Error sending PDF file: %v", err)
		http.Error(w, "Error sending PDF file", http.StatusInternalServerError)
	} else {
		log.Printf("PDF file for book ID %d successfully sent", book.ID)
	}
}
