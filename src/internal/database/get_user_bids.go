package database

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
	"github.com/google/uuid"
)

func (db Database) GetUserBids(ctx context.Context, getUserBidsDTO handlers.GetUserBidsDTO) (*[]handlers.BidOut, error) {
	query := `
		SELECT b.bid_id, b.bid_name, b.bid_description, b.bid_status, b.tender_id, b.bid_author_type, b.bid_author_id, b.bid_version, b.created_at
		FROM bid b
		JOIN employee e ON b.bid_author_id = e.id
		WHERE e.username = $1
		AND b.bid_version = (
			SELECT MAX(b2.bid_version)
			FROM bid b2
			WHERE b2.bid_id = b.bid_id
		)
		LIMIT $2 OFFSET $3
	`

	rows, err := db.db.QueryContext(ctx, query, getUserBidsDTO.Username, getUserBidsDTO.Limit, getUserBidsDTO.Offset)
	if err != nil {
		return nil, handlers.ErrMsgNotFound
	}
	defer rows.Close()

	var bids []handlers.BidOut

	for rows.Next() {
		var bid handlers.BidOut
		var bidId uuid.UUID
		var tenderId uuid.UUID
		var authorId uuid.UUID

		err := rows.Scan(&bidId, &bid.Name, &bid.Description, &bid.Status, &tenderId, &bid.AuthorType, &authorId, &bid.Version, &bid.CreatedAt)
		if err != nil {
			return nil, err
		}

		bid.Id = bidId.String()
		bid.TenderId = tenderId.String()
		bid.BidAuthorId = authorId.String()

		bids = append(bids, bid)
	}

	if err = rows.Err(); err != nil {
		return nil, handlers.ErrMsgNotFound
	}

	return &bids, nil
}
