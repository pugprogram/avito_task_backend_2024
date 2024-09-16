package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type GetBidsForTenderDTO struct {
	Limit    int64
	Offset   int64
	TenderId string
	Username string
}

func ToGetBidsForTenderDTO(tender GetBidsForTenderParams, tenderId TenderId) (*GetBidsForTenderDTO, string, error) {
	if tender.Limit == nil {
		var defaultLimit int32 = 5
		tender.Limit = &defaultLimit
	}

	if *tender.Limit < 0 && *tender.Limit > 50 {
		return nil, "limit", errors.New("invalid format for limit")
	}

	if tender.Offset == nil {
		var defaultOfsett int32 = 0
		tender.Offset = &defaultOfsett
	}

	if *tender.Offset < 0 {
		return nil, "offset", errors.New("invalid format for offset")
	}

	if tenderId == "" {
		return nil, "tenderId", errors.New("invalid format for id")
	}

	if tender.Username == "" || len(tender.Username) > 50 {
		return nil, "username", errors.New("invalid format for username")
	}

	return &GetBidsForTenderDTO{
		Limit:    int64(*tender.Limit),
		Offset:   int64(*tender.Offset),
		TenderId: string(tenderId),
		Username: string(tender.Username),
	}, "", nil

}

func (s Server) GetBidsForTender(w http.ResponseWriter, r *http.Request, tenderId TenderId, params GetBidsForTenderParams) {
	dto, _, err := ToGetBidsForTenderDTO(params, tenderId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		respJSON := ToErrorResponseJSON(err)
		_ = json.NewEncoder(w).Encode(respJSON)
		return
	}

	resp, err := s.repo.GetBidsForTender(r.Context(), *dto)
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

	respJSON := ToBidsJSON(*resp)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(respJSON)
}
