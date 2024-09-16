package database

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (db Database) GetBidReviews(ctx context.Context, getBidReviewsDTO handlers.GetBidReviewsDTO) (*[]handlers.BidReviewOut, error) {
	query := `
		SELECT bf.id, bf.created_at, bf.bid_feedback
		FROM bid_feedback bf
		JOIN bid b ON bf.bid_id = b.bid_id AND bf.bid_version = b.bid_version
		JOIN employee e ON e.id = b.bid_author_id
		WHERE e.username = $1
		AND b.bid_version = (
			SELECT MAX(bid_version)
			FROM bid
			WHERE bid_id = b.bid_id
		)
		ORDER BY bf.created_at
		LIMIT $2 OFFSET $3;

    `

	rows, err := db.db.QueryContext(ctx, query, getBidReviewsDTO.AuthorUsername, getBidReviewsDTO.Limit, getBidReviewsDTO.Offset)
	if err != nil {
		return nil, handlers.ErrMsgNotFound
	}
	defer rows.Close()

	var reviews []handlers.BidReviewOut
	for rows.Next() {
		var review handlers.BidReviewOut
		if err := rows.Scan(&review.Id, &review.CreatedAt, &review.Description); err != nil {
			return nil, handlers.ErrMsgNotFound
		}
		reviews = append(reviews, review)
	}

	if len(reviews) == 0 {
		return nil, handlers.ErrMsgNotFound
	}

	if err = rows.Err(); err != nil {
		return nil, handlers.ErrMsgNotFound
	}

	return &reviews, nil
}
