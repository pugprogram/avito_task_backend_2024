package repository

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (repo Repository) UpdateBidStatus(ctx context.Context, updateBidStatusDTO handlers.UpdateBidStatusDTO) (*handlers.BidOut, error) {
	err := repo.database.CheckUserExist(ctx, updateBidStatusDTO.Username)
	if err != nil {
		return nil, err
	}

	err = repo.database.CheckUserBidPermissionWithBidID(ctx, updateBidStatusDTO.Username, updateBidStatusDTO.BidID)
	if err != nil {
		return nil, err
	}

	updateBidStatusOut, err := repo.database.UpdateBidStatus(ctx, updateBidStatusDTO)
	if err != nil {
		return nil, err
	}

	return updateBidStatusOut, nil
}
