package repository

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (repo Repository) GetBidsForTender(ctx context.Context, getBidsForTenderDTO handlers.GetBidsForTenderDTO) (*[]handlers.BidOut, error) {
	// проверка что пользователь - работник
	err := repo.database.CheckUserExist(ctx, getBidsForTenderDTO.Username)
	if err != nil {
		return nil, err
	}

	// проверка что пользователь - ответственный за организацию
	err = repo.database.CheckUserPermissionWithTender(ctx, getBidsForTenderDTO.Username, getBidsForTenderDTO.TenderId)
	if err != nil {
		return nil, err
	}

	// получение предложений для тендера
	getBidsForTenderOut, err := repo.database.GetBidsForTender(ctx, getBidsForTenderDTO)
	if err != nil {
		return nil, err
	}

	return getBidsForTenderOut, nil
}
