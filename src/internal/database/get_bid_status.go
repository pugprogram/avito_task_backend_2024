package database

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
	"github.com/google/uuid"
)

func (db Database) GetBidStatus(ctx context.Context, getBidStatusDTO handlers.GetBidStatusDTO) (*string, error) {
	BidUUID, err := uuid.Parse(getBidStatusDTO.BidID)
	if err != nil {
		return nil, handlers.ErrMsgNotFound
	}

	query := `
        SELECT bid_status
        FROM bid
        WHERE bid_id = $1
        ORDER BY bid_version DESC
		LIMIT 1;
    `

	var status string
	err = db.db.QueryRowContext(ctx, query, BidUUID).Scan(&status)
	if err != nil {
		return nil, handlers.ErrMsgNotFound
	}

	return &status, nil
}
