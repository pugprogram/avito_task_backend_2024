package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type UpdateBidStatusDTO struct {
	BidID    string
	Status   string
	Username string
}

func ToUpdateBidStatusDTO(in UpdateBidStatusParams, bidId BidId) (*UpdateBidStatusDTO, string, error) {
	if len(bidId) > 100 {
		return nil, "bidId", errors.New("invalid format for bidId")
	}

	if in.Username == "" {
		return nil, "username", errors.New("invalid format for username")
	}

	status := in.Status

	flagStatusCanceled := status == BidStatusCanceled
	flagStatusPublished := status == BidStatusPublished

	if !(flagStatusCanceled || flagStatusPublished) {
		return nil, "status", errors.New("invalid format for status")
	}

	return &UpdateBidStatusDTO{
		BidID:    bidId,
		Username: in.Username,
		Status:   string(status),
	}, "", nil
}

func (s Server) UpdateBidStatus(w http.ResponseWriter, r *http.Request, bidId BidId, params UpdateBidStatusParams) {
	dto, _, err := ToUpdateBidStatusDTO(params, bidId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseJson := ToErrorResponseJSON(err)
		_ = json.NewEncoder(w).Encode(responseJson)
		return
	}

	resp, err := s.repo.UpdateBidStatus(r.Context(), *dto)
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

		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(&ErrorResponse{
			Reason: ErrMsgNotFound.Error(),
		})

		return
	}

	respJSON := ToBidJSON(*resp)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(respJSON)
}
