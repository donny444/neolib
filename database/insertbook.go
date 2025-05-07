package database

import (
	"context"
	"fmt"
)

func InsertBook(ctx context.Context, username string, isbn string, title string, publisher *string, category *string, author *string, pages *string, language *string, publicationYear *string, fileContent []byte) error {
	_, err := db.ExecContext(ctx, fmt.Sprintf(
		"INSERT INTO `%s` (isbn, title, publisher, category, author, pages, language, publication_year, file)"+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		username),
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
