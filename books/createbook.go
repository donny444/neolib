package books

import (
	"context"
	"fmt"
	"io"
	"log"
	"neolib/database"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

func CreateBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateBook function get called")

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
	fileExtension := ""

	// Retrieve the file from the form data
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		if err != http.ErrMissingFile {
			http.Error(w, "Error retrieving the file", http.StatusBadRequest)
			fmt.Println("Error: ", err)
			return
		}
	} else {
		if file == nil {
			fmt.Println("Client did not upload a file")
		} else {
			fmt.Println("File uploaded")
		}

		defer file.Close()

		if fileHeader != nil {
			fileExtension = filepath.Ext(fileHeader.Filename)
		}

		// Create a file in the images directory
		dst, err := os.Create(fmt.Sprintf("./images/%s%s", uuid, fileExtension))
		if err != nil {
			http.Error(w, "Unable to create the file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy the uploaded file to the destination file
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Unable to save the file", http.StatusInternalServerError)
			return
		}

		// Read the uploaded file content into a byte slice
		fileContent, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "Unable to read the file", http.StatusInternalServerError)
			fmt.Println("Error: ", err)
			return
		}
	}

	// Helper function to get pointer to string or nil
	optionalInput := func(value string) *string {
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

	// Continue with the database insertion
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use the InsertBook function from the database package
	err = database.InsertBook(ctx, uuid,
		requiredInput(r.FormValue("title")),
		requiredInput(r.FormValue("isbn")),
		optionalInput(r.FormValue("publisher")),
		optionalInput(r.FormValue("category")),
		optionalInput(r.FormValue("author")),
		optionalInput(r.FormValue("page")),
		optionalInput(r.FormValue("language")),
		optionalInput(r.FormValue("publication_year")),
		fileContent)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Book created"))
	result := fmt.Sprintf("Book created with UUID: %s", uuid)
	fmt.Println(result)
}
