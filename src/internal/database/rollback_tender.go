package database

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (db Database) RollbackTender(ctx context.Context, rollbackTenderDTO handlers.RollbackTenderDTO) (*handlers.TenderOUT, error) {

	queryGetLastVersion := `
	SELECT tender_id, tender_name, tender_description, service_type, status, organization_id, tender_version
	FROM tender
	WHERE tender_id = $1 AND tender_version = $2;
	`

	var tender handlers.TenderOUT
	var newTender handlers.TenderOUT

	err := db.db.QueryRowContext(ctx, queryGetLastVersion, rollbackTenderDTO.TenderID, rollbackTenderDTO.TenderVersion).Scan(
		&tender.Id,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.OrganizationId,
		&tender.Version,
	)
	if err != nil {
		return nil, handlers.ErrMsgNotFound
	}

	queryGetLastVersion = `
		SELECT MAX(tender_version)
		FROM tender
		WHERE tender_id = $1;
`

	var lastVersion int

	// Выполнение запроса для получения последней версии
	err = db.db.QueryRowContext(ctx, queryGetLastVersion, rollbackTenderDTO.TenderID).Scan(&lastVersion)
	if err != nil {
		return nil, handlers.ErrMsgNotFound
	}

	lastVersion = lastVersion + 1
	queryInsertNewVersion := `
		INSERT INTO tender (tender_id, tender_name, tender_description, service_type, status, organization_id, tender_version)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING tender_id, tender_name, tender_description, service_type, status, organization_id, tender_version, created_at;
	`

	err = db.db.QueryRowContext(ctx, queryInsertNewVersion, tender.Id, tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationId, lastVersion).Scan(
		&newTender.Id,
		&newTender.Name,
		&newTender.Description,
		&newTender.ServiceType,
		&newTender.Status,
		&newTender.OrganizationId,
		&newTender.Version,
		&newTender.CreatedAt,
	)
	if err != nil {
		return nil, handlers.ErrMsgNotFound
	}

	return &newTender, nil
}
