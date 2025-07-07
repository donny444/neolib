package database

import (
	"context"
	"database/sql"
	"fmt"
)

func SelectBooks(ctx context.Context, category *string, title *string, username string) (*sql.Rows, error) {
	if category != nil && title != nil {
		return db.QueryContext(ctx, fmt.Sprintf("SELECT isbn, title, file_extension FROM `%s_view` WHERE category = ? AND title LIKE ?", username), *category, "%"+*title+"%")
	} else if category != nil {
		return db.QueryContext(ctx, fmt.Sprintf("SELECT isbn, title, file_extension FROM `%s_view` WHERE category = ?", username), *category)
	} else if title != nil {
		return db.QueryContext(ctx, fmt.Sprintf("SELECT isbn, title, file_extension FROM `%s_view` WHERE title LIKE ?", username), "%"+*title+"%")
	} else {
		return db.QueryContext(ctx, fmt.Sprintf("SELECT isbn, title, file_extension FROM `%s_view`", username))
	}
}
