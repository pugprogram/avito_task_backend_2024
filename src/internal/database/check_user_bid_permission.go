package database

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (db Database) CheckUserBidPermission(ctx context.Context, username string) error {
	query := `
        SELECT EXISTS (
            SELECT 1
            FROM organization_employee
            WHERE username = $1
        )
    `

	var exists bool
	err := db.db.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return handlers.ErrMsgNotPermission
	}

	return nil
}
