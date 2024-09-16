package database

import (
	"context"
	"fmt"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (db Database) GetTenders(ctx context.Context, getTendersDTO handlers.GetTendersDTO) (*[]handlers.TenderOUT, error) {
	query := `
		SELECT t1.tender_id, t1.tender_name, t1.tender_description, t1.service_type, t1.status, t1.organization_id, t1.tender_version, t1.created_at
		FROM tender t1
		JOIN (
			SELECT tender_id, MAX(tender_version) AS max_version
			FROM tender
		`

	// Добавляем фильтр по типу услуги только если он задан
	if len(getTendersDTO.ServiceType) > 0 {
		query += fmt.Sprintf("WHERE service_type = '%s'", getTendersDTO.ServiceType[0])
	}

	for i := 1; i < len(getTendersDTO.ServiceType); i++ {
		query += fmt.Sprintf("OR service_type = '%s' ", getTendersDTO.ServiceType[i])
	}

	query += `
		GROUP BY tender_id) t2
		ON t1.tender_id = t2.tender_id AND t1.tender_version = t2.max_version
		LIMIT $1 OFFSET $2`

	rows, err := db.db.QueryContext(ctx, query, getTendersDTO.Limit, getTendersDTO.Offset)
	if err != nil {
		return nil, handlers.ErrNotValidParam
	}
	defer rows.Close()

	var tenders []handlers.TenderOUT

	for rows.Next() {
		var tender handlers.TenderOUT

		err := rows.Scan(&tender.Id, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationId, &tender.Version, &tender.CreatedAt)
		if err != nil {
			return nil, handlers.ErrNotValidParam
		}

		if tender.Status == "Published" {
			tenders = append(tenders, tender)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, handlers.ErrNotValidParam
	}

	if len(tenders) == 0 {
		return nil, handlers.ErrNotValidParam
	}

	return &tenders, nil
}
