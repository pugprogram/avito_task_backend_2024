package repository

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

type Database interface {
	CheckUserExist(context.Context, string) error
	CreateBid(context.Context, handlers.CreateBidDTO) (*handlers.BidOut, error)
	CreateTender(context.Context, handlers.CreateTenderDTO) (*handlers.TenderOUT, error)
	EditBid(context.Context, handlers.EditBidDTO) (*handlers.BidOut, error)
	EditTender(context.Context, handlers.EditTenderDTO) (*handlers.TenderOUT, error)
	GetBidReviews(context.Context, handlers.GetBidReviewsDTO) (*[]handlers.BidReviewOut, error)
	GetBidStatus(context.Context, handlers.GetBidStatusDTO) (*string, error)
	//Первая строка - Username, 2 - TenderId
	CheckUserPermissionWithTender(context.Context, string, string) error
	CheckUserPermissionWithOrganization(context.Context, string, string) error
	//Проверка на то, что человек или автор Bid или работает в организации Bid
	CheckUserBidPermissionWithBidID(context.Context, string, string) error
	//Пользователь работает в организации
	CheckUserExistID(context.Context, string) error
	CheckUserBidPermission(context.Context, string) error
	GetBidsForTender(context.Context, handlers.GetBidsForTenderDTO) (*[]handlers.BidOut, error)
	GetTenderStatus(context.Context, handlers.GetTenderStatusDTO) (*string, error)
	GetTenders(context.Context, handlers.GetTendersDTO) (*[]handlers.TenderOUT, error)
	//Пользователь из организации тендера имеет доступ к предложению
	CheckUserTenderPermissionToBidId(context.Context, string, string) error
	GetUserBids(context.Context, handlers.GetUserBidsDTO) (*[]handlers.BidOut, error)
	GetUserTenders(context.Context, handlers.GetUserTendersDTO) (*[]handlers.TenderOUT, error)
	RollbackBid(context.Context, handlers.RollbackBidDTO) (*handlers.BidOut, error)
	RollbackTender(context.Context, handlers.RollbackTenderDTO) (*handlers.TenderOUT, error)
	SubmitBidDecision(context.Context, handlers.SubmitBidDecisionDTO) (*handlers.BidOut, error)
	SubmitBidFeedback(context.Context, handlers.SubmitBidFeedbackDTO) (*handlers.BidOut, error)
	UpdateBidStatus(context.Context, handlers.UpdateBidStatusDTO) (*handlers.BidOut, error)
	UpdateTenderStatus(context.Context, handlers.UpdateTenderStatusDTO) (*handlers.TenderOUT, error)
}

type Repository struct {
	database Database
}

func NewRepository(database Database) Repository {
	return Repository{
		database: database,
	}
}
