package repository

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (repo Repository) UpdateTenderStatus(ctx context.Context, updateTenderStatusDTO handlers.UpdateTenderStatusDTO) (*handlers.TenderOUT, error) {
	// проверка что он работник
	err := repo.database.CheckUserExist(ctx, updateTenderStatusDTO.Username)
	if err != nil {
		return nil, err
	}

	// проверка что пользователь - ответственный за организацию
	err = repo.database.CheckUserPermissionWithTender(ctx, updateTenderStatusDTO.Username, updateTenderStatusDTO.TenderId)
	if err != nil {
		return nil, err
	}

	// обновление статуса тендера
	updateTenderStatusOUT, err := repo.database.UpdateTenderStatus(ctx, updateTenderStatusDTO)
	if err != nil {
		return nil, err
	}

	return updateTenderStatusOUT, nil
}
