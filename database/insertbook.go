package database

import "context"

func InsertBook(ctx context.Context, isbn string, title string, publisher *string, category *string, author *string, pages *string, language *string, publicationYear *string, fileContent []byte) error {
	_, err := db.ExecContext(ctx, "INSERT INTO books (isbn, title, publisher, category, author, pages, language, publication_year, file) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		isbn,
		title,
		publisher,
		category,
		author,
		pages,
		language,
		publicationYear,
		fileContent)
	return err
}
