package repository

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (repo Repository) GetTenderStatus(ctx context.Context, getTenderStatusDTO handlers.GetTenderStatusDTO) (*string, error) {
	// проверка что пользователь - работник
	err := repo.database.CheckUserExist(ctx, getTenderStatusDTO.Username)
	if err != nil {
		return nil, err
	}

	// проверка что пользователь - ответственный за организацию
	err = repo.database.CheckUserPermissionWithTender(ctx, getTenderStatusDTO.Username, getTenderStatusDTO.TenderId)
	if err != nil {
		return nil, err
	}

	// получение статуса тендера ответственного за организацию (может получить статус тендера только своей организации)
	getTenderStatusOut, err := repo.database.GetTenderStatus(ctx, getTenderStatusDTO)
	if err != nil {
		return nil, err
	}

	return getTenderStatusOut, nil
}
