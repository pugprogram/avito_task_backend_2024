package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type SubmitBidDecisionDTO struct {
	BidId    string
	Username string
	Decision string
}

func ToSubmitBidDecisionDTO(bidId BidId, params SubmitBidDecisionParams) (*SubmitBidDecisionDTO, string, error) {
	if len(bidId) > 100 {
		return nil, "bid id", errors.New("invalid format for bid id")
	}

	if params.Username == "" {
		return nil, "username", errors.New("invalid format for username")
	}

	flagDecigionReject := params.Decision == BidDecision(Rejected)
	flagDecigionApproved := params.Decision == BidDecision(Approved)

	if !flagDecigionReject && !flagDecigionApproved {
		return nil, "decision", errors.New("invalid format for decision")
	}

	return &SubmitBidDecisionDTO{
		BidId:    string(bidId),
		Username: params.Username,
		Decision: string(params.Decision),
	}, "", nil
}

func (s Server) SubmitBidDecision(w http.ResponseWriter, r *http.Request, bidId BidId, params SubmitBidDecisionParams) {
	dto, _, err := ToSubmitBidDecisionDTO(bidId, params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseJson := ToErrorResponseJSON(err)
		_ = json.NewEncoder(w).Encode(responseJson)
		return
	}

	resp, err := s.repo.SubmitBidDecision(r.Context(), *dto)
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
