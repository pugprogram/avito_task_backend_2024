package repository

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

func (repo Repository) EditTender(ctx context.Context, editTenderDto handlers.EditTenderDTO) (*handlers.TenderOUT, error) {
	// проверка что пользователь - работник
	err := repo.database.CheckUserExist(ctx, editTenderDto.Username)
	if err != nil {
		return nil, err
	}

	// проверка что пользователь - ответственный за организацию
	err = repo.database.CheckUserPermissionWithTender(ctx, editTenderDto.Username, editTenderDto.TenderId)
	if err != nil {
		return nil, err
	}

	// редактирование тендера
	editTenderOut, err := repo.database.EditTender(ctx, editTenderDto)
	if err != nil {
		return nil, err
	}

	return editTenderOut, nil
}
