package database

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
	"github.com/google/uuid"
)

func (db Database) CheckUserTenderPermissionToBidId(ctx context.Context, username string, bidId string) error {
	bidUUID, err := uuid.Parse(bidId)
	if err != nil {
		return err
	}

	var userID uuid.UUID
	err = db.db.QueryRowContext(ctx, `
		SELECT id
		FROM employee
		WHERE username = $1
	`, username).Scan(&userID)

	if err != nil {
		return err
	}

	var orgID uuid.UUID
	err = db.db.QueryRowContext(ctx, `
		SELECT t.organization_id
		FROM bid b
		JOIN tender t ON b.tender_id = t.tender_id
		WHERE b.bid_id = $1
	`, bidUUID).Scan(&orgID)

	if err != nil {
		return err
	}

	// Проверка, что пользователь отвечает за организацию
	var count int
	err = db.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM organization_responsible
		WHERE organization_id = $1 AND user_id = $2
	`, orgID, userID).Scan(&count)

	if err != nil {
		return err
	}

	// Проверяем, что пользователь ответственный за организацию
	if count == 0 {
		return handlers.ErrMsgNotFound
	}

	return nil
}
