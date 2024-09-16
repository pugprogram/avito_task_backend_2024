package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type CreateTenderDTO struct {
	CreatorUsername string
	Description     string
	Name            string
	OrganizationId  string
	ServiceType     string
	Status          string
}

func ToCreateTenderDTO(tender CreateTenderJSONBody) (*CreateTenderDTO, string, error) {
	if (len(tender.CreatorUsername) == 0) || (len(tender.CreatorUsername) > 50) {
		return nil, "creator username", errors.New("invalid format to creator username")
	}

	if len(tender.Description) == 0 || len(tender.Description) > 500 {
		return nil, "desription", errors.New("invalid format to description")
	}

	flagServiceTypeConstruction := tender.ServiceType != TenderServiceType(Construction)
	flagServiceTypeDelivery := tender.ServiceType != TenderServiceType(Delivery)
	flagServiceTypeManufacture := tender.ServiceType != TenderServiceType(Manufacture)

	if flagServiceTypeConstruction && flagServiceTypeDelivery && flagServiceTypeManufacture {
		return nil, "service type", errors.New("invalid format to service type")
	}

	if len(tender.OrganizationId) == 0 || len(tender.OrganizationId) > 100 {
		return nil, "Organization Id", errors.New("invalid format to organization id")
	}

	return &CreateTenderDTO{
		CreatorUsername: tender.CreatorUsername,
		Description:     tender.Description,
		Name:            tender.Name,
		OrganizationId:  tender.OrganizationId,
		ServiceType:     string(tender.ServiceType),
		Status:          "Created",
	}, "", nil
}

func (s Server) CreateTender(w http.ResponseWriter, r *http.Request) {
	var tenderReq CreateTenderJSONBody

	err := json.NewDecoder(r.Body).Decode(&tenderReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		respJSON := ToErrorResponseJSON(err)
		_ = json.NewEncoder(w).Encode(respJSON)

		return
	}

	dto, _, err := ToCreateTenderDTO(tenderReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		respJSON := ToErrorResponseJSON(err)
		_ = json.NewEncoder(w).Encode(respJSON)

		return
	}

	resp, err := s.repo.CreateTender(r.Context(), *dto)
	if err != nil {
		if errors.Is(err, ErrMsgUserNotExist) {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(&ErrorResponse{
				Reason: ErrMsgUserNotExist.Error(),
			})

			return
		} else {
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(&ErrorResponse{
				Reason: ErrMsgNotPermission.Error(),
			})

			return
		}
	}

	respJSON := ToTenderJSON(*resp)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(respJSON)
}
