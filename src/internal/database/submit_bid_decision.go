package database

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (db Database) SubmitBidDecision(ctx context.Context, submitBidDecision handlers.SubmitBidDecisionDTO) (*handlers.BidOut, error) {
	var (
		bidOut        handlers.BidOut
		bidOutPointer *handlers.BidOut
	)

	// Получение текущего статуса и количества одобрений
	query := `
		SELECT bid_status, approved_count, tender_id
		FROM bid
		WHERE bid_id = $1
		ORDER BY bid_version DESC
		LIMIT 1;
	`

	var (
		bidStatus     string
		approvedCount int
		tenderId      string
	)

	err := db.db.QueryRowContext(ctx, query, submitBidDecision.BidId).Scan(&bidStatus, &approvedCount, &tenderId)
	if err != nil {
		return nil, handlers.ErrMsgNotFound
	}

	if bidStatus == "Canceled" || bidStatus == "Rejected" || bidStatus == "Approved" {
		return nil, handlers.ErrMsgNotFound
	}

	if submitBidDecision.Decision == "Rejected" {
		bidOutPointer, err = db.UpdateBidStatus(ctx, handlers.UpdateBidStatusDTO{
			BidID:    submitBidDecision.BidId,
			Status:   "Rejected",
			Username: submitBidDecision.Username,
		})

		return bidOutPointer, err
	}

	// получаем количество ответственных за организацию
	var responsible_count int

	query = `
		SELECT COUNT(or_resp.user_id) AS responsible_count
		FROM organization_employee oe
		JOIN organization_responsible or_resp ON oe.organization_id = or_resp.organization_id
		WHERE oe.username = $1
		GROUP BY oe.organization_id;
	`

	err = db.db.QueryRowContext(ctx, query, submitBidDecision.Username).Scan(&responsible_count)
	if err != nil {
		return nil, handlers.ErrMsgNotFound
	}

	query = `
		SELECT approved_count 
		FROM bid 
		WHERE bid_id = $1 AND bid_version = 
			(SELECT max(bid_version)
			FROM bid
			WHERE bid_id = $1)
	`

	err = db.db.QueryRowContext(ctx, query, submitBidDecision.BidId).Scan(&approvedCount)
	if err != nil {
		return nil, handlers.ErrMsgNotFound
	}

	var kvorum int

	if 3 < responsible_count {
		kvorum = 3
	} else {
		kvorum = responsible_count
	}

	approvedCount += 1

	if approvedCount == kvorum {
		// устанавливаем предложению статус Approved
		bidOutPointer, err = db.UpdateBidStatus(ctx, handlers.UpdateBidStatusDTO{
			BidID:    submitBidDecision.BidId,
			Status:   "Approved",
			Username: submitBidDecision.Username,
		})
		if err != nil {
			return nil, handlers.ErrMsgNotFound
		}

		// закрываем тендер
		_, err = db.UpdateTenderStatus(ctx, handlers.UpdateTenderStatusDTO{
			TenderId: tenderId,
			Status:   "Closed",
			Username: submitBidDecision.Username,
		})
		if err != nil {
			return nil, handlers.ErrMsgNotFound
		}

		return bidOutPointer, nil
	}

	// обновляем количество аппрувов на предложении
	query = `
		UPDATE bid
		SET approved_count = $1
		WHERE bid_id = $2
			AND bid_version = (
			SELECT MAX(bid_version)
			FROM bid
			WHERE bid_id = $3
			)
		RETURNING bid_id, bid_name, bid_description, bid_status, tender_id, bid_author_type, bid_author_id, bid_version, created_at;
	`

	err = db.db.QueryRowContext(ctx, query, approvedCount, submitBidDecision.BidId, submitBidDecision.BidId).Scan(
		&bidOut.Id,
		&bidOut.Name,
		&bidOut.Description,
		&bidOut.Status,
		&bidOut.TenderId,
		&bidOut.AuthorType,
		&bidOut.BidAuthorId,
		&bidOut.Version,
		&bidOut.CreatedAt,
	)
	if err != nil {
		return nil, handlers.ErrMsgNotFound
	}

	return &bidOut, nil
}
