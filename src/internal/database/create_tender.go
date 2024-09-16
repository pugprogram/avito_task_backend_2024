package database

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
	"github.com/google/uuid"
)

func (db Database) CreateTender(ctx context.Context, createTenderDTO handlers.CreateTenderDTO) (*handlers.TenderOUT, error) {
	tenderID := uuid.New().String()

	query := `
        INSERT INTO tender (tender_id, tender_name, tender_description, 
        service_type, status, organization_id, tender_version)
        VALUES ($1, $2, $3, $4, $5, $6, 1)
        RETURNING tender_id, tender_name, tender_description, 
        service_type, status, organization_id, tender_version, created_at
    `

	var tenderOut handlers.TenderOUT
	err := db.db.QueryRowContext(ctx, query, tenderID, createTenderDTO.Name, createTenderDTO.Description,
		createTenderDTO.ServiceType, createTenderDTO.Status, createTenderDTO.OrganizationId).
		Scan(&tenderOut.Id, &tenderOut.Name, &tenderOut.Description,
			&tenderOut.ServiceType, &tenderOut.Status, &tenderOut.OrganizationId,
			&tenderOut.Version, &tenderOut.CreatedAt)

	if err != nil {
		return nil, handlers.ErrMsgNotPermission
	}

	return &tenderOut, nil
}
