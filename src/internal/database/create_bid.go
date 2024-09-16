package database

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
	"github.com/google/uuid"
)

func (db Database) CreateBid(ctx context.Context, createBidDTO handlers.CreateBidDTO) (*handlers.BidOut, error) {
	// получение статуса тендера
	var tenderStatus string

	bidID := uuid.New().String()
	query := `
		SELECT status
		FROM tender t
		WHERE t.tender_id = $1
		AND t.tender_version = (
			SELECT MAX(tender_version)
			FROM tender
			WHERE tender_id = $1
		)
		LIMIT 1;
	`

	err := db.db.QueryRowContext(ctx, query, createBidDTO.TenderId).Scan(&tenderStatus)
	if err != nil {
		return nil, handlers.ErrMsgNotPermission
	}

	// если тендер не опубликован, то создать предложение нельзя
	if tenderStatus != "Published" {
		return nil, handlers.ErrMsgNotPermission
	}

	queryMaxVersion := `
	SELECT MAX(tender_version)
	FROM tender
	WHERE tender_id = $1;
	`

	max_version := 0
	err = db.db.QueryRowContext(ctx, queryMaxVersion, createBidDTO.TenderId).Scan(&max_version)
	if err != nil {
		return nil, handlers.ErrMsgNotFound
	}

	// создание предложения
	query = `
		INSERT INTO bid (bid_id, bid_name, bid_description, 
		bid_status, tender_id, tender_version, approved_count, bid_author_type, bid_author_id, bid_version, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, 0, $7, $8, 1, CURRENT_TIMESTAMP)
		RETURNING bid_id, bid_name, bid_description, bid_status, tender_id, bid_author_type, bid_author_id, bid_version, created_at
	`

	var bidOut handlers.BidOut

	err = db.db.QueryRowContext(ctx, query, bidID, createBidDTO.Name, createBidDTO.Description,
		createBidDTO.Status, createBidDTO.TenderId, max_version, createBidDTO.AuthorType, createBidDTO.CreatorUsername).
		Scan(&bidOut.Id, &bidOut.Name, &bidOut.Description, &bidOut.Status,
			&bidOut.TenderId, &bidOut.AuthorType, &bidOut.BidAuthorId, &bidOut.Version, &bidOut.CreatedAt)
	if err != nil {
		return nil, handlers.ErrMsgNotFound
	}

	return &bidOut, nil
}
