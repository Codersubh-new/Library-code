package models

import (
	"github.com/subodh/library/pkg/config"
	"gorm.io/gorm"
)

var db *gorm.DB

type Book struct {
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Publication string `json:"publication"`
	// Path to the uploaded PDF
	PDFData []byte `json:"type:longblob"` //store pdf data as BLOB
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Book{})
}
func (b *Book) CreateBook() *Book {
	db.Create(&b)
	return b
}
func GetAllBook() []Book {
	var Book []Book
	db.Order("created_at desc").Find(&Book)
	return Book
}
func GetBookById(Id int64) (*Book, *gorm.DB) {
	var getBook Book
	db := db.Where("ID=?", Id).Find(&getBook)
	return &getBook, db
}
func DeleteBook(Id int64) Book {
	var book Book
	// Fetch the book by ID first
	db.Where("ID=?", Id).First(&book)
	// Pass a pointer to the book for deletion
	db.Delete(&book)
	return book
}
func SearchBooksByName(db *gorm.DB, name string) ([]Book, error) {
	var books []Book

	// Use GORM's Where method with LIKE for partial matching
	if err := db.Where("name LIKE ?", "%"+name+"%").Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}
