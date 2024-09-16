package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type RollbackBidDTO struct {
	BidId    string
	Username string
	Version  int32
}

func ToRollbackBidDTO(bidId BidId, version int32, params RollbackBidParams) (*RollbackBidDTO, string, error) {
	if len(bidId) > 100 {
		return nil, "bid id", errors.New("invalid format for bid id")
	}

	if params.Username == "" {
		return nil, "username", errors.New("invalid format for username")
	}

	if version < 1 {
		return nil, "version", errors.New("invalid format for version")
	}

	return &RollbackBidDTO{
		BidId:    string(bidId),
		Username: params.Username,
		Version:  version,
	}, "", nil
}

func (s Server) RollbackBid(w http.ResponseWriter, r *http.Request, bidId BidId, version int32, params RollbackBidParams) {
	dto, _, err := ToRollbackBidDTO(bidId, version, params)
	if err != nil {
		responseJson := ToErrorResponseJSON(err)

		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(responseJson)

		return
	}

	resp, err := s.repo.RollbackBid(r.Context(), *dto)
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
