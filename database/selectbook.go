package database

import (
	"context"
	"database/sql"
)

func SelectBook(ctx context.Context, username string, isbn string) (*sql.Row, error) {
	return db.QueryRowContext(ctx, "SELECT isbn, title, publisher, category, author, pages, language, publication_year FROM ? WHERE isbn = ?", username, isbn), nil
}
