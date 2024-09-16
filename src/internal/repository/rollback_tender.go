package repository

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (repo Repository) RollbackTender(ctx context.Context, rollbackTenderDTO handlers.RollbackTenderDTO) (*handlers.TenderOUT, error) {
	err := repo.database.CheckUserExist(ctx, rollbackTenderDTO.Username)
	if err != nil {
		return nil, err
	}

	err = repo.database.CheckUserPermissionWithTender(ctx, rollbackTenderDTO.Username, rollbackTenderDTO.TenderID)
	if err != nil {
		return nil, err
	}

	rollbackTenderOUT, err := repo.database.RollbackTender(ctx, rollbackTenderDTO)
	if err != nil {
		return nil, err
	}

	return rollbackTenderOUT, nil
}
