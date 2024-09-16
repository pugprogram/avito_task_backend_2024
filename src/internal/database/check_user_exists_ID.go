package database

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (db Database) CheckUserExistID(ctx context.Context, userName string) error {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM employee WHERE id = $1)`

	err := db.db.QueryRowContext(ctx, query, userName).Scan(&exists)
	if err != nil || !exists {
		return handlers.ErrMsgUserNotExist
	}

	return nil
}
