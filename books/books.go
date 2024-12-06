package books

import (
	"context"
	"io"
	"log"
	"neolib/database"
	"net/http"
	"text/template"
	"time"

	"github.com/google/uuid"
)

type Book struct {
	UUID            string
	Title           string
	Publisher       string
	Category        string
	Author          string
	Page            int
	Language        string
	PublicationYear int
	ISBN            string
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	// Generate a new UUID
	bookUUID := uuid.New().String()

	// Parse the multipart form, with a maximum memory of 10MB
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve the file from the form data
	file, _, err := r.FormFile("book_image")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// // Create the books directory if it doesn't exist
	// err = os.MkdirAll("books", os.ModePerm)
	// if err != nil {
	// 	http.Error(w, "Unable to create directory", http.StatusInternalServerError)
	// 	return
	// }

	// // Create a file in the books directory
	// dst, err := os.Create(fmt.Sprintf("books/%s", handler.Filename))
	// if err != nil {
	// 	http.Error(w, "Unable to create the file", http.StatusInternalServerError)
	// 	return
	// }
	// defer dst.Close()

	// // Copy the uploaded file to the destination file
	// _, err = io.Copy(dst, file)
	// if err != nil {
	// 	http.Error(w, "Unable to save the file", http.StatusInternalServerError)
	// 	return
	// }

	// Read the uploaded file content into a byte slice
	fileContent, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read the file", http.StatusInternalServerError)
		return
	}

	// Continue with the database insertion
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use the InsertBook function from the database package
	err = database.InsertBook(ctx, bookUUID, r.FormValue("title"), r.FormValue("publisher"), r.FormValue("category"), r.FormValue("author"), r.FormValue("page"), r.FormValue("language"), r.FormValue("publication_year"), r.FormValue("isbn"), fileContent)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusCreated)
}

func GetBooks(w http.ResponseWriter, _ *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := database.SelectBooks(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.UUID, &book.Title, &book.Publisher, &book.Category, &book.Author, &book.Page, &book.Language, &book.PublicationYear, &book.ISBN); err != nil {
			http.Error(w, "Unable to scan the row", http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	tmpl, err := template.ParseFiles("../templates/books.tmpl")
	if err != nil {
		http.Error(w, "Unable to load the template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, books); err != nil {
		http.Error(w, "Unable to execute the template", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func EditBook(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := database.UpdateBook(ctx, r.FormValue("title"), r.FormValue("publisher"), r.FormValue("category"), r.FormValue("author"), r.FormValue("page"), r.FormValue("language"), r.FormValue("publication_year"), r.FormValue("isbn"), r.FormValue("uuid"))
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := database.DeleteBook(ctx, r.URL.Query().Get("uuid"))
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusNoContent)
}
