package handler

import (
	"encoding/json"
	dto "modulegue/internal/delivery/mobile/customer/dto"
	"modulegue/internal/usecase/payment"
	"modulegue/pkg/response"
	"net/http"
)

type ScanHandler struct {
	submitQrUC   *payment.SubmitQrUseCase
	scanDetailUC *payment.GetScanDetailUseCase
}

func NewScanHandler(
	submitQrUC *payment.SubmitQrUseCase,
	scanDetailUC *payment.GetScanDetailUseCase,
) *ScanHandler {
	return &ScanHandler{submitQrUC: submitQrUC, scanDetailUC: scanDetailUC}
}

// Endpoint: POST /api/v2/linespot/scan
// Digunakan untuk: Submit QR dan langsung dapat detail (alur 1+2 dalam satu request)
func (h *ScanHandler) SubmitAndScan(w http.ResponseWriter, r *http.Request) {
	var req dto.SubmitQrRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	// Ambil customer ID dari JWT context
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Tahap 1: Submit QR → dapat session ID
	submitInput := payment.SubmitQrInput{
		ScannedQrString: req.ScannedQrString,
		CustomerID:      userID,
	}
	submitResult, err := h.submitQrUC.Execute(r.Context(), submitInput)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Tahap 2: Get Scan Detail menggunakan session ID
	detailInput := payment.GetScanDetailInput{
		SessionID:  submitResult.SessionID,
		CustomerID: userID,
	}
	detailResult, err := h.scanDetailUC.Execute(r.Context(), detailInput)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Mapping ke DTO
	resp := dto.ScanDetailDto{
		// CustomerId:   strconv.FormatInt(detailResult.CustomerID, 10),
		CustomerId:   detailResult.CustomerID,
		CustomerName: detailResult.CustomerName,
		Lokasi:       detailResult.LocationName,
		Duration:     detailResult.Duration,
		IsMember:     detailResult.IsMember,
		Total:        detailResult.Total,
		Breakdown:    []dto.PriceItemDto{},
	}
	for _, item := range detailResult.Breakdown {
		resp.Breakdown = append(resp.Breakdown, dto.PriceItemDto{
			Label:  item.Label,
			Amount: item.Amount,
		})
	}

	response.Success(w, http.StatusOK, "Scan berhasil", resp)
}
