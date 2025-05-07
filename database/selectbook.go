package database

import (
	"context"
	"database/sql"
	"fmt"
)

func SelectBook(ctx context.Context, username string, isbn string) (*sql.Row, error) {
	return db.QueryRowContext(ctx, fmt.Sprintf(
		"SELECT isbn, title, publisher, category, author, pages, language, publication_year FROM `%s`"+
			"WHERE isbn = ?", username), isbn), nil
}
