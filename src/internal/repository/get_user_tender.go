package repository

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (repo Repository) GetUserTenders(ctx context.Context, getUserTendersDTO handlers.GetUserTendersDTO) (*[]handlers.TenderOUT, error) {
	err := repo.database.CheckUserExist(ctx, getUserTendersDTO.Username)
	if err != nil {
		return nil, err
	}

	getUserTendersOUT, err := repo.database.GetUserTenders(ctx, getUserTendersDTO)
	if err != nil {
		return nil, err
	}

	return getUserTendersOUT, nil
}
