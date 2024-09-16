package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type GetBidReviewsDTO struct {
	TenderId          string
	AuthorUsername    string
	RequesterUsername string
	Limit             int64
	Offset            int64
}

func ToGetBidReviewsDTO(tenderId TenderId, params GetBidReviewsParams) (*GetBidReviewsDTO, string, error) {
	if len(tenderId) > 100 {
		return nil, "tender id", errors.New("invalid format for bid id")
	}

	if params.AuthorUsername == "" && len(params.AuthorUsername) > 50 {
		return nil, "authorUsername", errors.New("invalid format for authorUsername")
	}

	if params.RequesterUsername == "" || len(params.RequesterUsername) > 50 {
		return nil, "requesterUsername", errors.New("invalid format for requesterUsername")
	}

	if params.Limit == nil {
		var defaultLimit int32 = 5
		params.Limit = &defaultLimit
	}

	if *params.Limit < 0 && *params.Limit > 50 {
		return nil, "limit", errors.New("invalid format for limit")
	}

	if params.Offset == nil {
		var defaultOfsett int32 = 0
		params.Offset = &defaultOfsett
	}

	if *params.Offset < 0 {
		return nil, "offset", errors.New("invalid format for offset")
	}

	return &GetBidReviewsDTO{
		TenderId:          string(tenderId),
		AuthorUsername:    params.AuthorUsername,
		RequesterUsername: params.RequesterUsername,
		Limit:             int64(*params.Limit),
		Offset:            int64(*params.Offset),
	}, "", nil
}

func (s Server) GetBidReviews(w http.ResponseWriter, r *http.Request, tenderId TenderId, params GetBidReviewsParams) {
	dto, _, err := ToGetBidReviewsDTO(tenderId, params)
	if err != nil {
		respJSON := ToErrorResponseJSON(err)

		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(respJSON)

		return
	}

	resp, err := s.repo.GetBidReviews(r.Context(), *dto)
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

	respJSON := ToReviewsJSON(*resp)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(respJSON)
}
