package database

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (db Database) RollbackBid(ctx context.Context, rollbackBidDTO handlers.RollbackBidDTO) (*handlers.BidOut, error) {
	// Запрос для получения предложения по заданному ID и версии
	queryGetLastVersion := `
		SELECT bid_id, bid_name, bid_description, bid_status, tender_id, tender_version, approved_count, bid_author_type, bid_author_id, bid_version, created_at
		FROM bid
		WHERE bid_id = $1 AND bid_version = $2
		LIMIT 1;
	`

	var bid NewBid

	// Выполнение запроса для получения последней версии предложения
	err := db.db.QueryRowContext(ctx, queryGetLastVersion, rollbackBidDTO.BidId, rollbackBidDTO.Version).Scan(
		&bid.Id,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.TenderId,
		&bid.TenderVersion,
		&bid.ApprovedCount,
		&bid.AuthorType,
		&bid.BidAuthorId,
		&bid.Version,
		&bid.CreatedAt,
	)

	if err != nil {
		return nil, handlers.ErrMsgNotFound
	}

	// Получение последней версии предложения
	queryGetMaxVersion := `
		SELECT MAX(bid_version)
		FROM bid
		WHERE bid_id = $1;
	`

	var lastVersion int

	err = db.db.QueryRowContext(ctx, queryGetMaxVersion, rollbackBidDTO.BidId).Scan(&lastVersion)
	if err != nil {
		return nil, handlers.ErrMsgNotFound
	}

	// Инкрементирование версии
	lastVersion++

	// Вставляем новую версию предложения
	queryInsertNewVersion := `
		INSERT INTO bid (bid_id, bid_name, bid_description, bid_status, tender_id, tender_version, approved_count, bid_author_type, bid_author_id, bid_version, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())
		RETURNING bid_id, bid_name, bid_description, bid_status, tender_id, bid_author_type, bid_author_id, bid_version, created_at;
	`

	var newBid handlers.BidOut

	err = db.db.QueryRowContext(ctx, queryInsertNewVersion, bid.Id, bid.Name, bid.Description, bid.Status, bid.TenderId, bid.TenderVersion, bid.ApprovedCount, bid.AuthorType, bid.BidAuthorId, lastVersion).Scan(
		&newBid.Id,
		&newBid.Name,
		&newBid.Description,
		&newBid.Status,
		&newBid.TenderId,
		&newBid.AuthorType,
		&newBid.BidAuthorId,
		&newBid.Version,
		&newBid.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &newBid, nil
}
