package database

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
	"github.com/google/uuid"
)

// CheckUserPermissionWithTender проверяет что пользователь - ответственный за организацию
func (db Database) CheckUserPermissionWithTender(ctx context.Context, username string, tenderId string) error {
	tenderID, err := uuid.Parse(tenderId)
	if err != nil {
		return handlers.ErrMsgNotPermission
	}

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM organization_responsible org_responsible
			JOIN employee e ON org_responsible.user_id = e.id
			JOIN organization o ON org_responsible.organization_id = o.id
			JOIN tender t ON t.organization_id = o.id
			WHERE e.username = $1
			AND t.tender_id = $2
		);
	`

	var exists bool

	err = db.db.QueryRowContext(ctx, query, username, tenderID).Scan(&exists)
	if err != nil || !exists {
		return handlers.ErrMsgNotPermission
	}

	return nil
}
