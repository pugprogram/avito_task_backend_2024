package repository

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (repo Repository) GetTenders(ctx context.Context, getTendersDTO handlers.GetTendersDTO) (*[]handlers.TenderOUT, error) {
	getTendersOut, err := repo.database.GetTenders(ctx, getTendersDTO)
	if err != nil {
		return nil, err
	}

	return getTendersOut, nil
}
