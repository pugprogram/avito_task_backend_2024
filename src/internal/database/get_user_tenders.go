package database

import (
	"context"
	"database/sql"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (db Database) GetUserTenders(ctx context.Context, getUserTendersDTO handlers.GetUserTendersDTO) (*[]handlers.TenderOUT, error) {
	queryTenders := `
		SELECT t.*
		FROM tender t
		JOIN (
			SELECT tender_id, MAX(tender_version) AS max_version
			FROM tender
			GROUP BY tender_id
		) tm ON t.tender_id = tm.tender_id AND t.tender_version = tm.max_version
		JOIN organization_responsible orr ON t.organization_id = orr.organization_id
		JOIN employee e ON orr.user_id = e.id
		WHERE e.username = $1
		LIMIT $2 OFFSET $3
	`

	rows, err := db.db.QueryContext(ctx, queryTenders, getUserTendersDTO.Username, getUserTendersDTO.Limit, getUserTendersDTO.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenders []handlers.TenderOUT

	for rows.Next() {
		var tender handlers.TenderOUT

		err := rows.Scan(&tender.Id, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationId, &tender.Version, &tender.CreatedAt)
		if err != nil {
			return nil, handlers.ErrMsgNotFound
		}

		tenders = append(tenders, tender)
	}

	if err = rows.Err(); err != nil {
		return nil, handlers.ErrMsgNotFound
	}

	return &tenders, nil
}

func GetUserIDByUsername(ctx context.Context, db *sql.DB, Username string) (string, error) {
	var orgID string

	query := `
	SELECT id
	FROM employee
	WHERE username = $1
	LIMIT 1
`

	err := db.QueryRowContext(ctx, query, Username).Scan(&orgID)
	if err != nil {
		return "", err
	}

	return orgID, nil
}
