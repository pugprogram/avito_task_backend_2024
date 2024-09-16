package repository

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (repo Repository) GetBidReviews(ctx context.Context, getBidRewiewsDTO handlers.GetBidReviewsDTO) (*[]handlers.BidReviewOut, error) {
	// проверка что пользователь - работник
	err := repo.database.CheckUserExist(ctx, getBidRewiewsDTO.RequesterUsername)
	if err != nil {
		return nil, err
	}

	// проверка что запрашиваемый автор - работник
	err = repo.database.CheckUserExist(ctx, getBidRewiewsDTO.AuthorUsername)
	if err != nil {
		return nil, err
	}

	// проверка что пользователь - ответственный за организацию
	err = repo.database.CheckUserPermissionWithTender(ctx, getBidRewiewsDTO.RequesterUsername, getBidRewiewsDTO.TenderId)
	if err != nil {
		return nil, err
	}

	// получения отзывов на все предложения от запрашиваемого автора
	getBidReviewOut, err := repo.database.GetBidReviews(ctx, getBidRewiewsDTO)
	if err != nil {
		return nil, err
	}

	return getBidReviewOut, nil
}
