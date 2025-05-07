package database

import "context"

func DeleteBook(ctx context.Context, username string, isbn string) error {
	_, err := db.ExecContext(ctx, "DELETE FROM ? WHERE isbn = ?", username, isbn)
	return err
}
