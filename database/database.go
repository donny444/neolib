package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func SetupDatabase() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/neolib")
	if err != nil {
		log.Fatal(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}

func InsertBook(ctx context.Context, bookUUID, title, publisher, category, author, page, language, publicationYear, isbn string, fileContent []byte) error {
	_, err := db.ExecContext(ctx, "INSERT INTO books (uuid, title, publisher, category, author, page, language, publication_year, isbn, file_content) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		bookUUID,
		title,
		publisher,
		category,
		author,
		page,
		language,
		publicationYear,
		isbn,
		fileContent)
	return err
}

func SelectBooks(ctx context.Context) (*sql.Rows, error) {
	return db.QueryContext(ctx, "SELECT uuid, title, publisher, category, author, page, language, publication_year, isbn FROM books")
}

func UpdateBook(ctx context.Context, title, publisher, category, author, page, language, publicationYear, isbn, bookUUID string) error {
	_, err := db.ExecContext(ctx, "UPDATE books SET title = ?, publisher = ?, category = ?, author = ?, page = ?, language = ?, publication_year = ?, isbn = ? WHERE uuid = ?",
		title,
		publisher,
		category,
		author,
		page,
		language,
		publicationYear,
		isbn,
		bookUUID)
	return err
}

func DeleteBook(ctx context.Context, bookUUID string) error {
	_, err := db.ExecContext(ctx, "DELETE FROM books WHERE uuid = ?", bookUUID)
	return err
}
