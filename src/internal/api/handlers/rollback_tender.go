package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type RollbackTenderDTO struct {
	Username      string
	TenderID      string
	TenderVersion int32
}

func ToRollbackTenderDTO(tenderId TenderId, version int32, params RollbackTenderParams) (*RollbackTenderDTO, string, error) {
	if len(tenderId) > 100 {
		return nil, "tender ID", errors.New("invalid format for tender ID")
	}

	if version < 1 {
		return nil, "version tender", errors.New("invalid format for version tender")
	}

	return &RollbackTenderDTO{
		Username:      params.Username,
		TenderID:      tenderId,
		TenderVersion: version,
	}, "", nil
}

func (s Server) RollbackTender(w http.ResponseWriter, r *http.Request, tenderId TenderId, version int32, params RollbackTenderParams) {
	dto, _, err := ToRollbackTenderDTO(tenderId, version, params)
	if err != nil {
		responseJson := ToErrorResponseJSON(err)

		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(responseJson)

		return
	}

	resp, err := s.repo.RollbackTender(r.Context(), *dto)
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

	respJSOn := ToTenderJSON(*resp)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(respJSOn)
}
