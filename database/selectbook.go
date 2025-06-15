package database

import (
	"context"
	"database/sql"
	"fmt"
)

func SelectBook(ctx context.Context, username string, isbn string) (*sql.Row, error) {
	return db.QueryRowContext(ctx, fmt.Sprintf(
		"SELECT path, isbn, title, publisher, category, author, pages, language, publication_year, is_read FROM `%s_view`"+
			"WHERE isbn = ?", username), isbn), nil
}
