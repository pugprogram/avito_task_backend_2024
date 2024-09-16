package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type EditTenderDTO struct {
	TenderId          string
	Username          string
	TenderName        string
	TenderDescription string
	TenderServiceType string
}

func ToEditTenderDTO(tender EditTenderJSONBody, tenderId TenderId, params EditTenderParams) (*EditTenderDTO, string, error) {
	if len(tenderId) > 100 {
		return nil, "tender Id", errors.New("invalid format for tender id")
	}

	if tender.Name == nil {
		defaulTenderName := ""
		tender.Name = &defaulTenderName
	}

	if params.Username == "" || len(params.Username) > 50 {
		return nil, "username", errors.New("username not provided")
	}

	if len(*tender.Name) > 100 {
		return nil, "tender name", errors.New("invalid format for tender name")
	}

	if tender.Description == nil {
		defaulTenderDescription := ""
		tender.Description = &defaulTenderDescription
	}

	if len(*tender.Description) > 500 {
		return nil, "tender description", errors.New("invalid format for tender description")
	}

	if tender.ServiceType != nil {
		flagTenderConstruction := *tender.ServiceType != TenderServiceType(Construction)
		flagTenderDelivery := *tender.ServiceType != TenderServiceType(Delivery)
		flagTenderManufacture := *tender.ServiceType != TenderServiceType(Manufacture)
		if flagTenderConstruction && flagTenderDelivery && flagTenderManufacture {
			return nil, "tender service type", errors.New("invalid format for tender service type")
		}
	}

	if tender.ServiceType == nil {
		var defaultServiceType TenderServiceType = ""
		tender.ServiceType = &(defaultServiceType)
	}
	return &EditTenderDTO{
		TenderId:          tenderId,
		Username:          params.Username,
		TenderName:        *tender.Name,
		TenderDescription: *tender.Description,
		TenderServiceType: string(*tender.ServiceType),
	}, "", nil
}

func (s Server) EditTender(w http.ResponseWriter, r *http.Request, tenderId TenderId, params EditTenderParams) {
	var tenderReq EditTenderJSONBody
	err := json.NewDecoder(r.Body).Decode(&tenderReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		respJSON := ToErrorResponseJSON(err)
		_ = json.NewEncoder(w).Encode(respJSON)
		return
	}
	dto, _, err := ToEditTenderDTO(tenderReq, tenderId, params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		respJSON := ToErrorResponseJSON(err)
		_ = json.NewEncoder(w).Encode(respJSON)
		return
	}
	resp, err := s.repo.EditTender(r.Context(), *dto)
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

	respJSON := ToTenderJSON(*resp)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(respJSON)

}
