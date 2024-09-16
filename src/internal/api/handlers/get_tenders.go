package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type GetTendersDTO struct {
	Limit       int64
	Offset      int64
	ServiceType []TenderServiceType
}

func ToGetTendersDTO(tender GetTendersParams) (*GetTendersDTO, string, error) {
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

	if tender.ServiceType == nil {
		return &GetTendersDTO{
			Limit:       int64(*tender.Limit),
			Offset:      int64(*tender.Offset),
			ServiceType: []TenderServiceType{},
		}, "", nil
	}

	for _, serviceType := range *tender.ServiceType {
		flagServiceConstruction := serviceType != TenderServiceType(Construction)
		flagServiceDelivery := serviceType != TenderServiceType(Delivery)
		flagServiceManufacture := serviceType != TenderServiceType(Manufacture)
		if flagServiceConstruction && flagServiceDelivery && flagServiceManufacture {
			return nil, "service type", errors.New("invalid format for service type")
		}
	}

	return &GetTendersDTO{
		Limit:       int64(*tender.Limit),
		Offset:      (int64(*tender.Offset)),
		ServiceType: *tender.ServiceType,
	}, "", nil

}

func (s Server) GetTenders(w http.ResponseWriter, r *http.Request, params GetTendersParams) {
	dto, _, err := ToGetTendersDTO(params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseJson := ToErrorResponseJSON(err)
		_ = json.NewEncoder(w).Encode(responseJson)
		return
	}

	resp, err := s.repo.GetTenders(r.Context(), *dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseJson := ToErrorResponseJSON(err)
		_ = json.NewEncoder(w).Encode(responseJson)
		return
	}

	w.WriteHeader(http.StatusOK)

	respJSON := ToTendersResponseJSON(*resp)
	_ = json.NewEncoder(w).Encode(respJSON)
}
