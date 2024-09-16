package database

import (
	"context"
	"database/sql"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (db Database) CheckUserPermissionWithOrganization(ctx context.Context, username string, organizationID string) error {

	query := `
	SELECT EXISTS (
		SELECT 1
		FROM organization_responsible org
		WHERE org.user_id = (SELECT id FROM employee WHERE username = $1)
		AND org.organization_id = $2
)
    `

	var exists bool
	err := db.db.QueryRowContext(ctx, query, username, organizationID).Scan(&exists)
	if err != nil || !exists {
		return handlers.ErrMsgNotPermission
	}

	return nil
}

func GetOrganizationID(ctx context.Context, db *sql.DB, organizationName string) (string, error) {
	var orgID string
	query := `SELECT id FROM organization WHERE name = $1`
	err := db.QueryRowContext(ctx, query, organizationName).Scan(&orgID)
	if err != nil {
		return "", err
	}
	return orgID, nil
}
