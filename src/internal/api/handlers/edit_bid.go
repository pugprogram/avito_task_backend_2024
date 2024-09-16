package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type EditBidDTO struct {
	BidId       string
	Username    string
	BidName     string
	Description string
}

func ToEditBidDTO(in EditBidJSONBody, bidId BidId, params EditBidParams) (*EditBidDTO, string, error) {
	if len(bidId) > 100 {
		return nil, "bid id", errors.New("invalid format for bid id")
	}

	if in.Name == nil {
		defaulTenderName := ""
		in.Name = &defaulTenderName
	}

	if params.Username == "" || len(params.Username) > 50 {
		return nil, "username", errors.New("invalid format for username")
	}

	if len(*in.Name) > 100 {
		return nil, "bid name", errors.New("invalid format for bid name")
	}

	if in.Description == nil {
		defaulTenderDescription := ""
		in.Description = &defaulTenderDescription
	}

	if len(*in.Description) > 500 {
		return nil, "bid description", errors.New("invalid format for bid description")
	}

	return &EditBidDTO{
		BidId:       string(bidId),
		Username:    params.Username,
		BidName:     *in.Name,
		Description: *in.Description,
	}, "", nil
}

func (s Server) EditBid(w http.ResponseWriter, r *http.Request, bidId BidId, params EditBidParams) {
	var tenderReq EditBidJSONBody
	err := json.NewDecoder(r.Body).Decode(&tenderReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		respJSON := ToErrorResponseJSON(err)
		_ = json.NewEncoder(w).Encode(respJSON)
		return
	}

	dto, _, err := ToEditBidDTO(tenderReq, bidId, params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		respJSON := ToErrorResponseJSON(err)
		_ = json.NewEncoder(w).Encode(respJSON)
		return
	}

	resp, err := s.repo.EditBid(r.Context(), *dto)
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
