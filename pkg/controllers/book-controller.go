package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	models "github.com/subodh/library/pkg/models"
	utils "github.com/subodh/library/pkg/utils"
)

var NewBook models.Book

func GetBook(w http.ResponseWriter, r *http.Request) {
	newBook := models.GetAllBook()
	res, _ := json.Marshal(newBook)
	w.Header().Set("content.Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
func GetBookById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	ID, err := strconv.ParseInt(bookId, 0, 0)
	if err != nil {
		fmt.Println("error while parsing")
	}
	bookDetails, _ := models.GetBookById(ID)
	res, _ := json.Marshal(bookDetails)
	w.Header().Set("content.Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
func CreateBook(w http.ResponseWriter, r *http.Request) {
	CreateBook := &models.Book{}
	utils.ParseBody(r, CreateBook)
	b := CreateBook.CreateBook()
	res, _ := json.Marshal(b)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// Define the data structure to pass to the template
type BookDeleteData struct {
	Title     string
	BookStats []models.Book
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	BookID := vars["bookId"]
	ID, err := strconv.ParseInt(BookID, 0, 0)
	if err != nil {
		fmt.Println("error while parsing")
		http.Error(w, fmt.Sprintf("Invalid book ID: '%s'. Please provide a valid numeric ID.", BookID), http.StatusBadRequest)
		return
	}
	book := models.DeleteBook(ID)

	data := BookDeleteData{
		Title:     "BookDeletedData",
		BookStats: []models.Book{book},
	}
	// Pass data to template (after deletion)
	tmpl, err := template.ParseFiles("D:/practice/golang/devlelopment/newproject/library/views/home.html") // Assuming you have a template file
	if err != nil {
		fmt.Fprintf(w, "Error parsing template:%v", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	// Render template with book data
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}
}

// Define the data structure to pass to the template
type BookUpdateData struct {
	Title     string
	BookStats []models.Book
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	var updateBook = models.Book{}

	// Parse form data into updateBook struct (pass by reference)
	err := utils.ParseBody(r, &updateBook)
	if err != nil {
		fmt.Fprintf(w, "Error parsing form data: %v", err)
		return
	}
	fmt.Printf("Parsed Book Data: %+v\n", updateBook)

	vars := mux.Vars(r)
	bookId := vars["bookId"]
	ID, err := strconv.ParseInt(bookId, 0, 0)
	if err != nil {
		fmt.Fprintf(w, "Error parsing book ID: %v", err)
		return
	}

	// Fetch existing book details from DB
	bookDetails, db := models.GetBookById(ID)
	//fmt.Printf("Fetched Book Details: %+v\n", bookDetails)
	if updateBook.Name != "" {
		bookDetails.Name = updateBook.Name
	}
	if updateBook.Author != "" {
		bookDetails.Author = updateBook.Author
	}
	if updateBook.Publication != "" {
		bookDetails.Publication = updateBook.Publication
	}
	err = db.Save(&bookDetails).Error
	if err != nil {
		fmt.Printf("Error saving updated book: %v\n", err)
		return
	}

	// Prepare data for the template
	data := BookUpdateData{
		Title:     "BookUpdatedData",
		BookStats: []models.Book{*bookDetails},
	}
	// Pass data to template (after deletion)
	tmpl, err := template.ParseFiles("D:/practice/golang/devlelopment/newproject/library/views/home.html") // Assuming you have a template file
	if err != nil {
		fmt.Fprintf(w, "Error parsing template:%v", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	// Render template with book data
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}
}
