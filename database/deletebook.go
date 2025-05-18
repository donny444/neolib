package database

import (
	"context"
	"fmt"
)

func DeleteBook(ctx context.Context, username string, isbn string) error {
	_, err := db.ExecContext(ctx, fmt.Sprintf("DELETE FROM `%s` WHERE isbn = ?", username), isbn)
	return err
}
