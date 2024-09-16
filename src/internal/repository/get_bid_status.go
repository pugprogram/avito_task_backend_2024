package repository

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (repo Repository) GetBidStatus(ctx context.Context, getBidStatusDTO handlers.GetBidStatusDTO) (*string, error) {
	// проверка что пользователь - работник
	err := repo.database.CheckUserExist(ctx, getBidStatusDTO.Username)
	if err != nil {
		return nil, err
	}

	// проверка что пользователь - автор предложения или ответственный за организацию
	err = repo.database.CheckUserBidPermissionWithBidID(ctx, getBidStatusDTO.Username, getBidStatusDTO.BidID)
	if err != nil {
		return nil, err
	}

	// получить статус предложения
	bidStatus, err := repo.database.GetBidStatus(ctx, getBidStatusDTO)
	if err != nil {
		return nil, err
	}

	return bidStatus, nil
}
