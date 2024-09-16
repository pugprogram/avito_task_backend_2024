package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type SubmitBidFeedbackDTO struct {
	BidId       string
	Username    string
	BidFeedback string
}

func ToSubmitBidFeedbackDTO(bidId BidId, params SubmitBidFeedbackParams) (*SubmitBidFeedbackDTO, string, error) {
	if len(bidId) > 100 {
		return nil, "bid id", errors.New("invalid format for bid id")
	}

	if params.Username == "" {
		return nil, "username", errors.New("invalid format for username")
	}

	if len(params.BidFeedback) > 1000 {
		return nil, "bidFeedback", errors.New("invalid format for bidFeedback")
	}

	return &SubmitBidFeedbackDTO{
		BidId:       string(bidId),
		Username:    params.Username,
		BidFeedback: string(params.BidFeedback),
	}, "", nil
}

func (s Server) SubmitBidFeedback(w http.ResponseWriter, r *http.Request, bidId BidId, params SubmitBidFeedbackParams) {
	dto, _, err := ToSubmitBidFeedbackDTO(bidId, params)
	if err != nil {
		responseJson := ToErrorResponseJSON(err)

		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(responseJson)

		return
	}

	resp, err := s.repo.SubmitBidFeedback(r.Context(), *dto)
	if err != nil {
		if errors.Is(err, ErrMsgUserNotExist) {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(&ErrorResponse{
				Reason: ErrMsgUserNotExist.Error(),
			})

			return
		}

		if errors.Is(err, ErrMsgNotPermission) {
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(&ErrorResponse{
				Reason: ErrMsgNotPermission.Error(),
			})

			return
		}

		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(&ErrorResponse{
			Reason: ErrMsgNotFound.Error(),
		})

		return
	}

	respJSON := ToBidJSON(*resp)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(respJSON)
}
