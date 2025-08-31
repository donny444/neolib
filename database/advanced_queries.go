package database

import (
	"context"
	"database/sql"
	"fmt"
)

func SelectTopFive(ctx context.Context, username string) (*sql.Rows, error) {
	return db.QueryContext(ctx, fmt.Sprintf("SELECT category, COUNT(*) as count FROM `%s_view` GROUP BY category ORDER BY count DESC LIMIT 5", username))
}

func SelectCategoryStatuses(ctx context.Context, username string) (*sql.Rows, error) {
	return db.QueryContext(ctx, fmt.Sprintf("SELECT category, is_read, COUNT(*) as count FROM `%s_view` GROUP BY category, is_read", username))
}

/*
func SelectBooksByMonth(ctx context.Context, username string) (*sql.Rows, error) {
	return db.QueryContext(ctx, fmt.Sprintf("SELECT strftime('%%Y-%%m', start_date) as month, COUNT(*) as count FROM `%s_view` WHERE start_date IS NOT NULL AND is_read = 1 GROUP BY month ORDER BY month DESC", username))
}
*/
