package database

import (
	"context"
	"database/sql"
)

func SelectBook(ctx context.Context, isbn string) (*sql.Row, error) {
	return db.QueryRowContext(ctx, "SELECT isbn, title, publisher, category, author, pages, language, publication_year FROM books WHERE isbn = ?", isbn), nil
}
