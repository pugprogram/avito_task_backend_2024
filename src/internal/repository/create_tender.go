package repository

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (repo Repository) CreateTender(ctx context.Context, createTenderDTO handlers.CreateTenderDTO) (*handlers.TenderOUT, error) {
	// проверка что пользователь - работник
	err := repo.database.CheckUserExist(ctx, createTenderDTO.CreatorUsername)
	if err != nil {
		return nil, err
	}

	// проверка что пользователь - ответственный за организацию
	err = repo.database.CheckUserPermissionWithOrganization(ctx, createTenderDTO.CreatorUsername, createTenderDTO.OrganizationId)
	if err != nil {
		return nil, err
	}

	// создание тендера
	createTenderOut, err := repo.database.CreateTender(ctx, createTenderDTO)
	if err != nil {
		return nil, err
	}

	return createTenderOut, nil
}
