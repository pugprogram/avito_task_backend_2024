package database

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
	"github.com/google/uuid"
)

func (db Database) EditTender(ctx context.Context, editTenderDTO handlers.EditTenderDTO) (*handlers.TenderOUT, error) {
	tenderUUID, err := uuid.Parse(editTenderDTO.TenderId)
	if err != nil {
		return nil, handlers.ErrMsgNotFound
	}

	queryGetLastVersion := `
		SELECT tender_id, tender_name, tender_description, service_type, status, organization_id, tender_version
		FROM tender
		WHERE tender_id = $1
		ORDER BY tender_version DESC
		LIMIT 1;
	 `

	var tender handlers.TenderOUT
	var newTender handlers.TenderOUT

	// Выполнение запроса для получения последней версии тендера
	err = db.db.QueryRowContext(ctx, queryGetLastVersion, tenderUUID).Scan(
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

	// создвем новую версию тендера
	if len(editTenderDTO.TenderDescription) != 0 {
		tender.Description = editTenderDTO.TenderDescription
	}

	if len(editTenderDTO.TenderServiceType) != 0 {
		tender.ServiceType = editTenderDTO.TenderServiceType
	}
	if len(editTenderDTO.TenderName) != 0 {
		tender.Name = editTenderDTO.TenderName
	}

	tender.Version++ // Инкрементирование версии

	// вставляем новую версию тендера
	queryInsertNewVersion := `
		INSERT INTO tender (tender_id, tender_name, tender_description, service_type, status, organization_id, tender_version)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING tender_id, tender_name, tender_description, service_type, status, organization_id, tender_version, created_at;
	`

	err = db.db.QueryRowContext(ctx, queryInsertNewVersion, tender.Id, tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationId, tender.Version).Scan(
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
