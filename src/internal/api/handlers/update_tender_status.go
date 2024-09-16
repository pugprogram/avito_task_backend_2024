package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type UpdateTenderStatusDTO struct {
	TenderId string
	Status   string
	Username string
}

func ToUpdateTenderStatusDTO(tenderId TenderId, params UpdateTenderStatusParams) (*UpdateTenderStatusDTO, string, error) {
	if len(tenderId) > 100 {
		return nil, "tender ID", errors.New("invalid format for tender ID")
	}

	var flagStatusPublished = params.Status != TenderStatus(TenderStatusPublished)
	var flagStatusClosed = params.Status != TenderStatus(TenderStatusClosed)

	if flagStatusPublished && flagStatusClosed {
		return nil, "status", errors.New("invalid format for status")
	}
	return &UpdateTenderStatusDTO{
		TenderId: tenderId,
		Status:   string(params.Status),
		Username: string(params.Username),
	}, "", nil

}

func (s Server) UpdateTenderStatus(w http.ResponseWriter, r *http.Request, tenderId TenderId, params UpdateTenderStatusParams) {
	dto, _, err := ToUpdateTenderStatusDTO(tenderId, params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseJson := ToErrorResponseJSON(err)
		_ = json.NewEncoder(w).Encode(responseJson)
		return
	}

	resp, err := s.repo.UpdateTenderStatus(r.Context(), *dto)
	if err != nil {
		if errors.Is(err, ErrMsgUserNotExist) {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(&ErrorResponse{
				Reason: ErrMsgNotFound.Error(),
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

	respJSON := ToTenderJSON(*resp)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(respJSON)
}
