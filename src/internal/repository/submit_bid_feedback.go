package repository

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (repo Repository) SubmitBidFeedback(ctx context.Context, submitBidFeedbackDTO handlers.SubmitBidFeedbackDTO) (*handlers.BidOut, error) {
	err := repo.database.CheckUserExist(ctx, submitBidFeedbackDTO.Username)
	if err != nil {
		return nil, err
	}

	err = repo.database.CheckUserBidPermissionWithBidID(ctx, submitBidFeedbackDTO.Username, submitBidFeedbackDTO.BidId)
	if err != nil {
		return nil, err
	}

	submitBidFeedbackOUT, err := repo.database.SubmitBidFeedback(ctx, submitBidFeedbackDTO)
	if err != nil {
		return nil, err
	}

	return submitBidFeedbackOUT, nil
}
