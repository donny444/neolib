package database

import (
	"context"
	"fmt"
	"log"
)

func InsertUser(ctx context.Context, username string, email string, password string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal("Error starting transaction: ", err)
	}

	// Insert user into the users table
	_, err = tx.ExecContext(ctx, "INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, password)
	if err != nil {
		tx.Rollback()
		log.Fatal("Error inserting user: ", err)
	}

	// Create a table for the user
	_, err = tx.ExecContext(ctx, fmt.Sprintf("CREATE TABLE `%s` ( "+
		"`title` varchar(255) NOT NULL, "+
		"`publisher` varchar(255) DEFAULT NULL, "+
		"`category` varchar(255) DEFAULT NULL, "+
		"`author` varchar(255) DEFAULT NULL, "+
		"`pages` smallint(5) UNSIGNED DEFAULT NULL, "+
		"`language` varchar(255) DEFAULT NULL, "+
		"`publication_year` year(4) DEFAULT NULL, "+
		"`isbn` varchar(20) NOT NULL, "+
		"`file` blob DEFAULT NULL, "+
		"`created_at` timestamp NOT NULL DEFAULT current_timestamp(), "+
		"`updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()) "+
		"ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;", username))
	if err != nil {
		tx.Rollback()
		log.Fatal("Error creating user table: ", err)
	}

	// Create a view for the user
	_, err = tx.ExecContext(ctx, fmt.Sprintf("CREATE VIEW `%s_view` AS SELECT * FROM `%s`", username, username))
	if err != nil {
		tx.Rollback()
		log.Fatal("Error creating user view: ", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		log.Fatal("Error committing transaction: ", err)
	}

	return nil
}
