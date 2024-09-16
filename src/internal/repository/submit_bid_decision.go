package repository

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (repo Repository) SubmitBidDecision(ctx context.Context, submitBidDecisionDTO handlers.SubmitBidDecisionDTO) (*handlers.BidOut, error) {
	err := repo.database.CheckUserExist(ctx, submitBidDecisionDTO.Username)
	if err != nil {
		return nil, err
	}

	err = repo.database.CheckUserTenderPermissionToBidId(ctx, submitBidDecisionDTO.Username, submitBidDecisionDTO.BidId)
	if err != nil {
		return nil, err
	}

	SubmitBidDecisionOUT, err := repo.database.SubmitBidDecision(ctx, submitBidDecisionDTO)
	if err != nil {
		return nil, err
	}

	return SubmitBidDecisionOUT, nil
}
