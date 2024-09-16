package database

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
	"github.com/google/uuid"
)

func (db Database) SubmitBidFeedback(ctx context.Context, submitBidFeedbackDTO handlers.SubmitBidFeedbackDTO) (*handlers.BidOut, error) {
	// Проверка на существование employee с данным username
	employeeQuery := `
		SELECT id FROM employee WHERE username = $1
	`

	var employeeID uuid.UUID

	err := db.db.QueryRowContext(ctx, employeeQuery, submitBidFeedbackDTO.Username).Scan(&employeeID)
	if err != nil {
		return nil, err
	}

	// Получение последней версии предложения
	bidVersionQuery := `
		SELECT bid_version
		FROM bid
		WHERE bid_id = $1
		ORDER BY bid_version DESC
		LIMIT 1
	`

	var bidVersion int

	err = db.db.QueryRowContext(ctx, bidVersionQuery, submitBidFeedbackDTO.BidId).Scan(&bidVersion)
	if err != nil {
		return nil, err
	}

	// Вставка отзыва в таблицу bid_feedback
	insertQuery := `
		INSERT INTO bid_feedback (id, bid_id, bid_version, author, bid_feedback, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`

	feedbackID := uuid.New() // Генерация уникального идентификатора для отзыва

	_, err = db.db.ExecContext(ctx, insertQuery, feedbackID, submitBidFeedbackDTO.BidId, bidVersion, employeeID, submitBidFeedbackDTO.BidFeedback)
	if err != nil {
		return nil, err
	}

	// Получение информации о предложении после добавления отзыва
	bid, err := db.GetBid(ctx, submitBidFeedbackDTO.BidId)
	if err != nil {
		return nil, err
	}

	return bid, nil
}

func (db Database) GetBid(ctx context.Context, bidID string) (*handlers.BidOut, error) {
	query := `
		SELECT bid_id, bid_name, bid_description, bid_status, tender_id, bid_author_type, bid_author_id, bid_version, created_at
		FROM bid
		WHERE bid_id = $1
		ORDER BY bid_version DESC
		LIMIT 1
	`

	var (
		bidOut      handlers.BidOut
		bidAuthorID uuid.UUID
		tenderID    uuid.UUID
	)

	err := db.db.QueryRowContext(ctx, query, bidID).Scan(
		&bidOut.Id,
		&bidOut.Name,
		&bidOut.Description,
		&bidOut.Status,
		&tenderID,
		&bidOut.AuthorType,
		&bidAuthorID,
		&bidOut.Version,
		&bidOut.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	bidOut.TenderId = tenderID.String()
	bidOut.BidAuthorId = bidAuthorID.String()

	return &bidOut, nil
}
