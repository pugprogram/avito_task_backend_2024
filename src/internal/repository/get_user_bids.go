package repository

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (repo Repository) GetUserBids(ctx context.Context, getUserBidsDTO handlers.GetUserBidsDTO) (*[]handlers.BidOut, error) {
	err := repo.database.CheckUserExist(ctx, getUserBidsDTO.Username)
	if err != nil {
		return nil, err
	}

	getUserBidsOUT, err := repo.database.GetUserBids(ctx, getUserBidsDTO)
	if err != nil {
		return nil, err
	}

	return getUserBidsOUT, nil
}
