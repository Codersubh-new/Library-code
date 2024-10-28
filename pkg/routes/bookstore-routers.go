package routes

import (
	"github.com/gorilla/mux"
	controllers "github.com/subodh/library/pkg/controllers"
)

var RegisterBookStoreRoutes = func(router *mux.Router) {
	router.HandleFunc("/book/", controllers.CreateBook).Methods("POST")
	router.HandleFunc("/book/", controllers.GetBook).Methods("GET")
	router.HandleFunc("/book/{bookId}", controllers.GetBookById).Methods("GET")
	router.HandleFunc("/book/{bookId}", controllers.UpdateBook).Methods("POST")
	router.HandleFunc("/book/delete/{bookId}", controllers.DeleteBook).Methods("POST")
	router.HandleFunc("/report", controllers.ReportHandler).Methods("GET")
	router.HandleFunc("/report/csv", controllers.ReportCSVHandler).Methods("GET")
	router.HandleFunc("/upload/csv", controllers.UploadCSVHandler).Methods("POST")
	router.HandleFunc("/books/search", controllers.SearchBookByName).Methods("GET")
	router.HandleFunc("/upload/pdf", controllers.UploadPDFHandler).Methods("POST")
	router.HandleFunc("/download/{bookId}", controllers.DownloadPDFHandler).Methods("GET")
}
