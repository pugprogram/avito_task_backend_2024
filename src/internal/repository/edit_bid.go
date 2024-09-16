package repository

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (repo Repository) EditBid(ctx context.Context, editBidDTO handlers.EditBidDTO) (*handlers.BidOut, error) {
	// проверка что пользователь - работник
	err := repo.database.CheckUserExist(ctx, editBidDTO.Username)
	if err != nil {
		return nil, err
	}

	// проверка что пользователь - автор предложения или ответственный за организацию
	err = repo.database.CheckUserBidPermissionWithBidID(ctx, editBidDTO.Username, editBidDTO.BidId)
	if err != nil {
		return nil, err
	}

	// редактирование предложения
	editBidOut, err := repo.database.EditBid(ctx, editBidDTO)
	if err != nil {
		return nil, err
	}

	return editBidOut, nil
}
