package database

import "context"

func DeleteBook(ctx context.Context, uuid string) error {
	_, err := db.ExecContext(ctx, "DELETE FROM books WHERE uuid = ?", uuid)
	return err
}
