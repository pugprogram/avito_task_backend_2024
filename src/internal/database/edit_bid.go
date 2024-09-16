package database

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

type NewBid struct {
	Id            string
	Name          string
	Description   string
	Status        string
	TenderId      string
	TenderVersion int32
	ApprovedCount int32
	AuthorType    string
	BidAuthorId   string
	Version       int32
	CreatedAt     string
}

func (db Database) EditBid(ctx context.Context, editBid handlers.EditBidDTO) (*handlers.BidOut, error) {
	queryGetLastVersion := `
		SELECT bid_id, bid_name, bid_description, bid_status, tender_id, tender_version, approved_count, bid_author_type, bid_author_id, bid_version, created_at
		FROM bid
		WHERE bid_id = $1
		ORDER BY bid_version DESC
		LIMIT 1;
	 `

	var bid handlers.BidOut

	var newBid NewBid

	// Выполнение запроса для получения последней версии предложения
	err := db.db.QueryRowContext(ctx, queryGetLastVersion, editBid.BidId).Scan(
		&newBid.Id,
		&newBid.Name,
		&newBid.Description,
		&newBid.Status,
		&newBid.TenderId,
		&newBid.TenderVersion,
		&newBid.ApprovedCount,
		&newBid.AuthorType,
		&newBid.BidAuthorId,
		&newBid.Version,
		&newBid.CreatedAt,
	)

	if err != nil {
		return nil, handlers.ErrMsgNotFound
	}

	// создвем новую версию предложения
	if len(editBid.Description) != 0 {
		newBid.Description = editBid.Description
	}

	if len(editBid.BidName) != 0 {
		newBid.Name = editBid.BidName
	}

	newBid.Version++ // Инкрементирование версии

	// вставляем новую версию предложения
	queryInsertNewVersion := `
		INSERT INTO bid (bid_id, bid_name, bid_description, bid_status, tender_id, tender_version, approved_count, bid_author_type, bid_author_id, bid_version, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, 0, $7, $8, $9, NOW())
		RETURNING bid_id, bid_name, bid_description, bid_status, tender_id, bid_author_type, bid_author_id, bid_version, created_at;
	`

	err = db.db.QueryRowContext(ctx, queryInsertNewVersion, newBid.Id, newBid.Name, newBid.Description, newBid.Status, newBid.TenderId, newBid.TenderVersion, newBid.AuthorType, newBid.BidAuthorId, newBid.Version).Scan(
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
		return nil, handlers.ErrMsgNotFound
	}

	return &bid, nil
}
