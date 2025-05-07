package database

import "context"

func UpdateBook(ctx context.Context, username string, isbn string, title string, publisher *string, category *string, author *string, pages *string, language *string, publicationYear *string) error {
	_, err := db.ExecContext(ctx, "UPDATE ? SET title = ?, publisher = ?, category = ?, author = ?, pages = ?, language = ?, publication_year = ? WHERE isbn = ?",
		username,
		title,
		publisher,
		category,
		author,
		pages,
		language,
		publicationYear,
		isbn,
	)
	return err
}
