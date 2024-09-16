package handlers

import (
	"context"
)

type Repository interface {
	GetUserBids(context.Context, GetUserBidsDTO) (*[]BidOut, error)
	CreateBid(context.Context, CreateBidDTO) (*BidOut, error)
	EditBid(context.Context, EditBidDTO) (*BidOut, error)
	SubmitBidFeedback(context.Context, SubmitBidFeedbackDTO) (*BidOut, error)
	RollbackBid(context.Context, RollbackBidDTO) (*BidOut, error)
	GetBidStatus(context.Context, GetBidStatusDTO) (*string, error)
	UpdateBidStatus(context.Context, UpdateBidStatusDTO) (*BidOut, error)
	SubmitBidDecision(context.Context, SubmitBidDecisionDTO) (*BidOut, error)
	GetBidsForTender(context.Context, GetBidsForTenderDTO) (*[]BidOut, error)
	GetBidReviews(context.Context, GetBidReviewsDTO) (*[]BidReviewOut, error)
	CheckServer()
	GetTenders(context.Context, GetTendersDTO) (*[]TenderOUT, error)
	GetUserTenders(context.Context, GetUserTendersDTO) (*[]TenderOUT, error)
	CreateTender(context.Context, CreateTenderDTO) (*TenderOUT, error)
	EditTender(context.Context, EditTenderDTO) (*TenderOUT, error)
	RollbackTender(context.Context, RollbackTenderDTO) (*TenderOUT, error)
	GetTenderStatus(context.Context, GetTenderStatusDTO) (*string, error)
	UpdateTenderStatus(context.Context, UpdateTenderStatusDTO) (*TenderOUT, error)
}

type Server struct {
	repo Repository
}

func NewServer(repository Repository) Server {
	return Server{
		repo: repository,
	}
}
