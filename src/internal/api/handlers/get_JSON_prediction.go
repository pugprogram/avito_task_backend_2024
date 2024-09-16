package handlers

func ToErrorResponseJSON(errorResponse error) ErrorResponse {
	errorResponseJson := ErrorResponse{
		Reason: errorResponse.Error(),
	}

	return errorResponseJson
}

func ToTendersResponseJSON(tenders []TenderOUT) []Tender {
	var tendersResponseJSON []Tender

	var tenderResponseJSON Tender

	for _, tender := range tenders {
		tenderResponseJSON = Tender{
			CreatedAt:      tender.CreatedAt,
			Description:    tender.Description,
			Id:             tender.Id,
			Name:           tender.Name,
			OrganizationId: tender.OrganizationId,
			ServiceType:    TenderServiceType(tender.ServiceType),
			Status:         TenderStatus(tender.Status),
			Version:        int32(tender.Version),
		}
		tendersResponseJSON = append(tendersResponseJSON, tenderResponseJSON)
	}

	return tendersResponseJSON
}

func ToBidJSON(bid BidOut) Bid {
	return Bid{
		AuthorId:    bid.BidAuthorId,
		AuthorType:  BidAuthorType(bid.AuthorType),
		CreatedAt:   bid.CreatedAt,
		Description: bid.Description,
		Id:          bid.Id,
		Name:        bid.Name,
		Status:      BidStatus(bid.Status),
		TenderId:    bid.TenderId,
		Version:     bid.Version,
	}
}

func ToBidsJSON(bids []BidOut) []Bid {
	var bidsResponseJSON []Bid

	var bidResponseJSON Bid

	for _, bid := range bids {
		bidResponseJSON = ToBidJSON(bid)
		bidsResponseJSON = append(bidsResponseJSON, bidResponseJSON)
	}

	return bidsResponseJSON
}

func ToTenderJSON(tender TenderOUT) Tender {
	return Tender{
		CreatedAt:      tender.CreatedAt,
		Description:    tender.Description,
		Id:             tender.Id,
		Name:           tender.Name,
		OrganizationId: tender.OrganizationId,
		ServiceType:    TenderServiceType(tender.ServiceType),
		Status:         TenderStatus(tender.Status),
		Version:        int32(tender.Version),
	}
}

func ToReviewsJSON(bidReviews []BidReviewOut) []BidReview {
	var (
		bidsResponseJSON []BidReview
		bidResponseJSON  BidReview
	)

	for _, bid := range bidReviews {
		bidResponseJSON = BidReview{
			CreatedAt:   bid.CreatedAt,
			Description: bid.Description,
			Id:          bid.Id,
		}
		bidsResponseJSON = append(bidsResponseJSON, bidResponseJSON)
	}

	return bidsResponseJSON
}
