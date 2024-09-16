package repository

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (repo Repository) CreateBid(ctx context.Context, createBidInput handlers.CreateBidDTO) (*handlers.BidOut, error) {
	// проверка что пользователь - работник
	err := repo.database.CheckUserExistID(ctx, createBidInput.CreatorUsername)
	if err != nil {
		return nil, err
	}

	// создание предложения в таблице
	createBidOut, err := repo.database.CreateBid(ctx, createBidInput)
	if err != nil {
		return nil, err
	}

	return createBidOut, nil
}
