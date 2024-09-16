package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type GetBidStatusDTO struct {
	BidID    string
	Username string
}

func ToGetBidStatusDTO(in GetBidStatusParams, bidId BidId) (*GetBidStatusDTO, string, error) {
	if len(bidId) > 100 {
		return nil, "bidId", errors.New("invalid format for bidId")
	}

	if in.Username == "" || len(in.Username) > 50 {
		return nil, "username", errors.New("invalid format for username")
	}

	return &GetBidStatusDTO{
		BidID:    bidId,
		Username: in.Username,
	}, "", nil
}

func (s Server) GetBidStatus(w http.ResponseWriter, r *http.Request, bidId BidId, params GetBidStatusParams) {
	dto, _, err := ToGetBidStatusDTO(params, bidId)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(&ErrorResponse{
			Reason: ErrMsgUserNotExist.Error(),
		})

		return
	}

	resp, err := s.repo.GetBidStatus(r.Context(), *dto)
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

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]BidStatus{
		"status": BidStatus(*resp),
	})
}
