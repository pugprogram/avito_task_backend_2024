package database

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
	"github.com/google/uuid"
)

func (db Database) UpdateBidStatus(ctx context.Context, in handlers.UpdateBidStatusDTO) (*handlers.BidOut, error) {
	bidID, err := uuid.Parse(in.BidID)
	if err != nil {
		return nil, err
	}

	query := `
		UPDATE bid
		SET bid_status = $1
		WHERE bid_id = $2
		AND bid_version = (
			SELECT MAX(bid_version) 
			FROM bid 
			WHERE bid_id = $3
		)
		RETURNING bid_id, bid_name, bid_description, bid_status, tender_id, bid_author_type, bid_author_id, bid_version, created_at
	`

	var bid handlers.BidOut

	err = db.db.QueryRowContext(ctx, query, in.Status, bidID, bidID).Scan(
		&bid.Id,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.TenderId,
		&bid.AuthorType,
		&bid.BidAuthorId,
		&bid.Version,
		&bid.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &bid, nil
}
