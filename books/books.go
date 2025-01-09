package books

import (
	"context"
	"fmt"
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
	Publisher       *string
	Category        *string
	Author          *string
	Page            *int
	Language        *string
	PublicationYear *int
	ISBN            string
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	// Generate a new UUID
	uuid := uuid.New().String()
	fmt.Println("UUID: ", uuid)

	// Parse the multipart form, with a maximum memory of 10MB
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		fmt.Println("Error: ", err)
		return
	}

	var fileContent []byte

	// Retrieve the file from the form data
	file, _, err := r.FormFile("book_image")
	if err != nil {
		if err != http.ErrMissingFile {
			http.Error(w, "Error retrieving the file", http.StatusBadRequest)
			fmt.Println("Error: ", err)
			return
		}
	} else {
		defer file.Close()

		// Read the uploaded file content into a byte slice
		fileContent, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "Unable to read the file", http.StatusInternalServerError)
			fmt.Println("Error: ", err)
			return
		}
	}

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

	// Helper function to get pointer to string or nil
	getStringPointer := func(value string) *string {
		if value == "" {
			return nil
		}
		return &value
	}

	requiredInput := func(value string) string {
		if value == "" {
			http.Error(w, "Required inputs are not filled", http.StatusBadRequest)
		}
		return value
	}

	// Helper function to get pointer to int or nil
	// getIntPointer := func(value string) *int {
	// 	if value == "" {
	// 		return nil
	// 	}
	// 	intValue, err := strconv.Atoi(value)
	// 	if err != nil {
	// 		return nil
	// 	}
	// 	return &intValue
	// }

	// Continue with the database insertion
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use the InsertBook function from the database package
	err = database.InsertBook(ctx, uuid,
		requiredInput(r.FormValue("title")),
		requiredInput(r.FormValue("isbn")),
		getStringPointer(r.FormValue("publisher")),
		getStringPointer(r.FormValue("category")),
		getStringPointer(r.FormValue("author")),
		getStringPointer(r.FormValue("page")),
		getStringPointer(r.FormValue("language")),
		getStringPointer(r.FormValue("publication_year")),
		fileContent)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Book created"))
	result := fmt.Sprintf("Book created with UUID: %s", uuid)
	fmt.Println(result)
}

func GetBooks(w http.ResponseWriter, _ *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := database.SelectBooks(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		http.Error(w, "Unable to get the columns", http.StatusInternalServerError)
		return
	}
	fmt.Println("Columns: ", columns)

	var books []Book
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.UUID, &book.Title, &book.ISBN); err != nil {
			http.Error(w, "Unable to scan the row", http.StatusInternalServerError)
			fmt.Println("Scan error: ", err)
			return
		}
		books = append(books, book)
	}

	tmpl, err := template.ParseFiles("templates/books.tmpl")
	if err != nil {
		http.Error(w, "Unable to load the template", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to load the template"))
		return
	}

	if err := tmpl.Execute(w, books); err != nil {
		http.Error(w, "Unable to execute the template", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to execute the template"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Retrieved all books"))
	fmt.Println("Retrieved all books in the system")
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")
	fmt.Println("UUID: ", uuid)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row, err := database.SelectBook(ctx, uuid)
	if err != nil {
		log.Fatal(err)
	}

	var book Book
	if err := row.Scan(&book.Title, &book.Publisher, &book.Category, &book.Author, &book.Page, &book.Language, &book.PublicationYear, &book.ISBN); err != nil {
		http.Error(w, "Unable to scan the row", http.StatusInternalServerError)
		fmt.Println("Scan error: ", err)
		return
	}

	tmpl, err := template.ParseFiles("templates/book.tmpl")
	if err != nil {
		http.Error(w, "Unable to load the template", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to load the template"))
		return
	}

	if err := tmpl.Execute(w, book); err != nil {
		http.Error(w, "Unable to execute the template", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to execute the template"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Retrieved the book specified by UUID"))
	result := fmt.Sprintf("Retrieved the book specified by the UUID: %s", uuid)
	fmt.Println(result)
}

func EditBook(w http.ResponseWriter, r *http.Request) {
	uuid := r.PathValue("book")
	fmt.Println("UUID: ", uuid)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Helper function to get pointer to string or nil
	getStringPointer := func(value string) *string {
		if value == "" {
			return nil
		}
		return &value
	}

	requiredInput := func(value string) string {
		if value == "" {
			http.Error(w, "Required inputs are not filled", http.StatusBadRequest)
		}
		return value
	}

	err := database.UpdateBook(ctx, uuid,
		requiredInput(r.FormValue("title")),
		requiredInput(r.FormValue("isbn")),
		getStringPointer(r.FormValue("publisher")),
		getStringPointer(r.FormValue("category")),
		getStringPointer(r.FormValue("author")),
		getStringPointer(r.FormValue("page")),
		getStringPointer(r.FormValue("language")),
		getStringPointer(r.FormValue("publication_year")))
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to update the book"))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Book updated"))
	result := fmt.Sprintf("Book with UUID: %s is updated", uuid)
	fmt.Println(result)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	uuid := r.PathValue("book")
	fmt.Println("UUID: ", uuid)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := database.DeleteBook(ctx, uuid)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to delete the book"))
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("The book specified by UUID is deleted"))
	result := fmt.Sprintf("Deleted the book specified by the UUID: %s", uuid)
	fmt.Println(result)
}
