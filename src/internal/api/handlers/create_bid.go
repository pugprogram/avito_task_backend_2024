package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type CreateBidDTO struct {
	CreatorUsername string
	Description     string
	Name            string
	AuthorType      string
	Status          string
	TenderId        string
}

func ToCreateBidDTO(bid CreateBidJSONBody) (*CreateBidDTO, string, error) {
	if len(bid.AuthorId) > 100 {
		return nil, "creator name", errors.New("invalid format for creator name")
	}

	if len(bid.Description) > 500 {
		return nil, "description", errors.New("invalid format for description")
	}

	if len(bid.TenderId) > 100 {
		return nil, "id", errors.New("invalid format for id")
	}

	if len(bid.Name) > 100 {
		return nil, "name", errors.New("invalid format for name")
	}

	flagUserOrganization := bid.AuthorType != BidAuthorType(Organization)
	flagUser := bid.AuthorType != BidAuthorType(User)

	if flagUserOrganization && flagUser {
		return nil, "status", errors.New("invalid format for author type")
	}

	return &CreateBidDTO{
		CreatorUsername: bid.AuthorId,
		Description:     bid.Description,
		Name:            bid.Name,
		AuthorType:      string(bid.AuthorType),
		Status:          "Created",
		TenderId:        bid.TenderId,
	}, "", nil
}

func (s Server) CreateBid(w http.ResponseWriter, r *http.Request) {
	var tenderReq CreateBidJSONBody

	err := json.NewDecoder(r.Body).Decode(&tenderReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		respJSON := ToErrorResponseJSON(err)
		_ = json.NewEncoder(w).Encode(respJSON)

		return
	}

	bid, _, err := ToCreateBidDTO(tenderReq)
	if err != nil {
		respJSON := ToErrorResponseJSON(err)

		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(respJSON)

		return
	}

	resp, err := s.repo.CreateBid(r.Context(), *bid)
	if err != nil {
		if errors.Is(err, ErrMsgUserNotExist) {
			w.WriteHeader(http.StatusUnauthorized)

			respJSON := ToErrorResponseJSON(err)
			_ = json.NewEncoder(w).Encode(respJSON)

			return
		}

		if errors.Is(err, ErrMsgNotPermission) {
			w.WriteHeader(http.StatusForbidden)

			respJSON := ToErrorResponseJSON(err)
			_ = json.NewEncoder(w).Encode(respJSON)

			return
		}

		w.WriteHeader(http.StatusNotFound)

		respJSON := ToErrorResponseJSON(err)
		_ = json.NewEncoder(w).Encode(respJSON)

		return
	}

	respJSON := ToBidJSON(*resp)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(respJSON)
}
