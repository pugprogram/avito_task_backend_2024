package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type GetTenderStatusDTO struct {
	Username string
	TenderId string
}

func ToGetTenderStatusDTO(tenderId TenderId, params GetTenderStatusParams) (*GetTenderStatusDTO, string, error) {
	if len(tenderId) > 100 {
		return nil, "tender ID", errors.New("invalid format for tender ID")
	}

	if params.Username == nil {
		var defaultUsername = ""
		params.Username = &defaultUsername
	}

	return &GetTenderStatusDTO{
		Username: *params.Username,
		TenderId: tenderId,
	}, "", nil
}

func (s Server) GetTenderStatus(w http.ResponseWriter, r *http.Request, tenderId TenderId, params GetTenderStatusParams) {
	dto, _, err := ToGetTenderStatusDTO(tenderId, params)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		respJSON := ToErrorResponseJSON(err)
		_ = json.NewEncoder(w).Encode(respJSON)
		return
	}

	tenderStatus, err := s.repo.GetTenderStatus(r.Context(), *dto)
	if err != nil {
		if errors.Is(err, ErrMsgUserNotExist) {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(&ErrorResponse{
				Reason: ErrMsgUserNotExist.Error(),
			})
			return
		}

		if errors.Is(err, ErrMsgNotPermission) {
			w.WriteHeader(http.StatusUnauthorized)
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
	_ = json.NewEncoder(w).Encode(map[string]TenderStatus{
		"status": TenderStatus(*tenderStatus),
	})

}
