package database

import "context"

func UpdateBook(ctx context.Context, isbn string, title string, publisher *string, category *string, author *string, page *string, language *string, publicationYear *string) error {
	_, err := db.ExecContext(ctx, "UPDATE books SET title = ?, publisher = ?, category = ?, author = ?, page = ?, language = ?, publication_year = ? WHERE isbn = ?",
		title,
		publisher,
		category,
		author,
		page,
		language,
		publicationYear,
		isbn,
	)
	return err
}
