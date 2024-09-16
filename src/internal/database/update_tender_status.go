package database

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (db Database) UpdateTenderStatus(ctx context.Context, in handlers.UpdateTenderStatusDTO) (*handlers.TenderOUT, error) {
	if in.Status == "Rejected" {
		return nil, handlers.ErrMsgNotPermission
	}

	query := `
		UPDATE tender
		SET status = $1
		WHERE tender_id = $2
		AND tender_version = (
			SELECT MAX(tender_version) 
			FROM tender 
			WHERE tender_id = $2
		)
		RETURNING tender_id, tender_name, tender_description, service_type, status, organization_id, tender_version, created_at;
	`

	var tender handlers.TenderOUT

	err := db.db.QueryRowContext(ctx, query, in.Status, in.TenderId).Scan(
		&tender.Id,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.OrganizationId,
		&tender.Version,
		&tender.CreatedAt,
	)
	if err != nil {
		return nil, handlers.ErrMsgNotFound
	}

	return &tender, nil
}
