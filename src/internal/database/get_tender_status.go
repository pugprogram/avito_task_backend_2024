package database

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (db Database) GetTenderStatus(ctx context.Context, getTenderStatusDTO handlers.GetTenderStatusDTO) (*string, error) {

	query := `
		SELECT t.status
		FROM tender t
		WHERE t.tender_id = $1 
		AND tender_version = (
		SELECT MAX(tender_version) 
		FROM tender 
		WHERE tender_id = $1
		);
	`

	var status string
	err := db.db.QueryRowContext(ctx, query, getTenderStatusDTO.TenderId).Scan(&status)
	if err != nil {
		return nil, handlers.ErrMsgNotFound
	}

	return &status, nil
}
