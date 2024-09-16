package repository

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (repo Repository) RollbackBid(ctx context.Context, rollbackBidDTO handlers.RollbackBidDTO) (*handlers.BidOut, error) {
	err := repo.database.CheckUserExist(ctx, rollbackBidDTO.Username)
	if err != nil {
		return nil, err
	}

	err = repo.database.CheckUserBidPermissionWithBidID(ctx, rollbackBidDTO.Username, rollbackBidDTO.BidId)
	if err != nil {
		return nil, err
	}

	rollbackBidOUT, err := repo.database.RollbackBid(ctx, rollbackBidDTO)
	if err != nil {
		return nil, err
	}

	return rollbackBidOUT, nil
}
