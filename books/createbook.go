package books

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"neolib/database"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func CreateBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateBook function get called")

	username, ok := r.Context().Value("username").(string)
	if !ok || username == "" {
		log.Fatal("Username not found in context")
	}

	isbn := r.FormValue("isbn")

	// Parse the multipart form, with a maximum memory of 10MB
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		fmt.Println("Error: ", err)
		return
	}

	// Retrieve the file from the form data
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		if err != http.ErrMissingFile {
			http.Error(w, "Error retrieving the file", http.StatusBadRequest)
			fmt.Println("Error: ", err)
			return
		}
	}

	fileContent := createFile(username, isbn, file, fileHeader)

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

	// Declare and assign the image path if a file was uploaded, otherwise set it to nil
	var imagePath *string
	if fileContent != nil {
		path := fmt.Sprintf("/images/%s/%s", username, isbn)
		imagePath = &path
	} else {
		imagePath = nil
	}

	err = database.InsertBook(ctx,
		username,
		requiredInput(r.FormValue("isbn")),
		requiredInput(r.FormValue("title")),
		optionalInput(r.FormValue("publisher")),
		optionalInput(r.FormValue("category")),
		optionalInput(r.FormValue("author")),
		optionalInput(r.FormValue("pages")),
		optionalInput(r.FormValue("language")),
		optionalInput(r.FormValue("publication_year")),
		fileContent,
		imagePath,
	)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Book created"))
	result := fmt.Sprintf("Book created with ISBN: %s", isbn)
	fmt.Println(result)
}

func createFile(username, isbn string, file multipart.File, fileHeader *multipart.FileHeader) []byte {
	if file == nil {
		fmt.Println("Client did not upload a file")
		return nil
	}
	fmt.Println("File uploaded")
	defer file.Close()

	fileExtension := ""
	if fileHeader != nil {
		fileExtension = filepath.Ext(fileHeader.Filename)
	}

	// Create a file in the images directory
	dst, err := os.Create(fmt.Sprintf("./images/%s/%s%s", username, isbn, fileExtension))
	if err != nil {
		log.Fatal("Unable to create the file: ", err)
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	_, err = io.Copy(dst, file)
	if err != nil {
		log.Fatal("Unable to save the file: ", err)
	}

	// Read the uploaded file content into a byte slice
	fileContent, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Unable to read the file: ", err)
	}

	return fileContent
}
