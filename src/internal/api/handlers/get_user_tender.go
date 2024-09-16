package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type GetUserTendersDTO struct {
	Limit    int64
	Offset   int64
	Username string
}

func ToGetUserTenderDTO(tender GetUserTendersParams) (*GetUserTendersDTO, string, error) {
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

	if tender.Username == nil || *tender.Username == "" {
		return nil, "username", errors.New("invalid format for username")
	}

	return &GetUserTendersDTO{
		Limit:    int64(*tender.Limit),
		Offset:   (int64(*tender.Offset)),
		Username: *tender.Username,
	}, "", nil
}

func (s Server) GetUserTenders(w http.ResponseWriter, r *http.Request, params GetUserTendersParams) {
	dto, _, err := ToGetUserTenderDTO(params)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(&ErrorResponse{
			Reason: ErrMsgUserNotExist.Error(),
		})

		return
	}

	resp, err := s.repo.GetUserTenders(r.Context(), *dto)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(&ErrorResponse{
			Reason: ErrMsgUserNotExist.Error(),
		})

		return
	}

	respJSON := ToTendersResponseJSON(*resp)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(respJSON)
}
