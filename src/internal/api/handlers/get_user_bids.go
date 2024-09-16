package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type GetUserBidsDTO struct {
	Limit    int64
	Offset   int64
	Username string
}

func toGetUserBidsDTO(in GetUserBidsParams) (*GetUserBidsDTO, string, error) {
	if in.Limit == nil {
		var defaultLimit int32 = 5
		in.Limit = &defaultLimit
	}

	if *in.Limit < 0 && *in.Limit > 50 {
		return nil, "limit", errors.New("invalid format for limit")
	}

	if in.Offset == nil {
		var defaultOfsett int32 = 0
		in.Offset = &defaultOfsett
	}

	if *in.Offset < 0 {
		return nil, "offset", errors.New("invalid format for offset")
	}

	if in.Username == nil || *in.Username == "" {
		return nil, "username", errors.New("invalid format for username")
	}

	return &GetUserBidsDTO{
		Limit:    int64(*in.Limit),
		Offset:   (int64(*in.Offset)),
		Username: *in.Username,
	}, "", nil
}

func (s Server) GetUserBids(w http.ResponseWriter, r *http.Request, params GetUserBidsParams) {
	dto, _, err := toGetUserBidsDTO(params)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(&ErrorResponse{
			Reason: ErrMsgUserNotExist.Error(),
		})
		return
	}

	resp, err := s.repo.GetUserBids(r.Context(), *dto)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(&ErrorResponse{
			Reason: ErrMsgUserNotExist.Error(),
		})

		return
	}

	respJSON := ToBidsJSON(*resp)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(respJSON)
}
