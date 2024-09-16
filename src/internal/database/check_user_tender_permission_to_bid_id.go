package database

import (
	"context"
	"database/sql"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
	"github.com/google/uuid"
)

func (db Database) CheckUserTenderPermission(ctx context.Context, username string, tenderID string) error {
	tenderUUID, err := uuid.Parse(tenderID)
	if err != nil {
		return handlers.ErrMsgNotPermission
	}

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM tender t
			JOIN organization_responsible o_r ON t.organization_id = o_r.organization_id
			JOIN employee e ON o_r.user_id = e.id
			WHERE t.tender_id = $1
			  AND e.username = $2
		)
	`

	var exists bool
	err = db.db.QueryRowContext(ctx, query, tenderUUID, username).Scan(&exists)

	if err != nil {
		return err
	}

	if !exists {
		return sql.ErrNoRows
	}

	return nil
}
