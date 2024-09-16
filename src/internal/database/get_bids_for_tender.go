package database

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (db Database) GetBidsForTender(ctx context.Context, getBidsForTender handlers.GetBidsForTenderDTO) (*[]handlers.BidOut, error) {
	query := `
		SELECT b.bid_id, b.bid_name, b.bid_description, b.bid_status, b.tender_id, b.bid_author_type, b.bid_author_id, b.bid_version, b.created_at
		FROM bid b
		JOIN (
			SELECT bid_id, MAX(bid_version) AS max_version
			FROM bid
			WHERE bid_status = 'Published'
			GROUP BY bid_id
		) AS latest_bid ON b.bid_id = latest_bid.bid_id AND b.bid_version = latest_bid.max_version
		WHERE b.bid_status = 'Published' AND b.tender_id = $1
		ORDER BY b.bid_name
		LIMIT $2 OFFSET $3;
	`

	rows, err := db.db.QueryContext(ctx, query, getBidsForTender.TenderId, getBidsForTender.Limit, getBidsForTender.Offset)
	if err != nil {
		return nil, handlers.ErrMsgNotFound
	}
	defer rows.Close()

	var bids []handlers.BidOut

	for rows.Next() {
		var bid handlers.BidOut

		err := rows.Scan(&bid.Id, &bid.Name, &bid.Description, &bid.Status, &bid.TenderId, &bid.AuthorType, &bid.BidAuthorId, &bid.Version, &bid.CreatedAt)
		if err != nil {
			return nil, handlers.ErrMsgNotFound
		}

		bids = append(bids, bid)
	}

	if err = rows.Err(); err != nil {
		return nil, handlers.ErrMsgNotFound
	}

	return &bids, nil
}
