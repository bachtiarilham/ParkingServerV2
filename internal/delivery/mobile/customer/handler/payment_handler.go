package handler

import (
	"encoding/json"
	dto "modulegue/internal/delivery/mobile/customer/dto"
	"modulegue/internal/usecase/payment"
	"modulegue/pkg/response"
	"net/http"
	"strconv"
)

type TransactionHandler struct {
	submitQrUC   *payment.SubmitQrUseCase
	scanDetailUC *payment.GetScanDetailUseCase
	paymentUC    *payment.ExecutePaymentUseCase
}

func NewTransactionHandler(
	submitQrUC *payment.SubmitQrUseCase,
	scanDetailUC *payment.GetScanDetailUseCase,
	paymentUC *payment.ExecutePaymentUseCase,
) *TransactionHandler {
	return &TransactionHandler{
		submitQrUC:   submitQrUC,
		scanDetailUC: scanDetailUC,
		paymentUC:    paymentUC,
	}
}

func (h *TransactionHandler) SubmitQr(w http.ResponseWriter, r *http.Request) {
	var req dto.SubmitQrRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Ambil customer ID dari context JWT
	customerID, ok := r.Context().Value("userID").(int64)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	input := payment.SubmitQrInput{
		ScannedQrString: req.ScannedQrString,
		CustomerID:      customerID,
	}

	result, err := h.submitQrUC.Execute(r.Context(), input)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Kita return session ID untuk digunakan di tahap selanjutnya
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    result.Message,
		"session_id": strconv.FormatInt(result.SessionID, 10), // Kirim session ID ke client
	})
}

func (h *TransactionHandler) GetScanDetail(w http.ResponseWriter, r *http.Request) {
	sessionIDStr := r.URL.Query().Get("session_id") // Ambil dari query param
	if sessionIDStr == "" {
		response.Error(w, http.StatusBadRequest, "session_id is required")
		return
	}
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid session_id")
		return
	}

	customerID, ok := r.Context().Value("userID").(int64)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	input := payment.GetScanDetailInput{
		SessionID:  sessionID,
		CustomerID: customerID,
	}

	result, err := h.scanDetailUC.Execute(r.Context(), input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := dto.ScanDetailDto{
		CustomerId:   result.CustomerID,
		CustomerName: result.CustomerName,
		Lokasi:       result.LocationName,
		Duration:     result.Duration,
		IsMember:     result.IsMember,
		Total:        result.Total,
		Breakdown:    []dto.PriceItemDto{},
	}
	for _, item := range result.Breakdown {
		resp.Breakdown = append(resp.Breakdown, dto.PriceItemDto{
			Label:  item.Label,
			Amount: item.Amount,
		})
	}

	response.Success(w, http.StatusOK, "Scan detail retrieved", resp)
}

func (h *TransactionHandler) ExecutePayment(w http.ResponseWriter, r *http.Request) {
	var req dto.ExecutePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	customerID, ok := r.Context().Value("userID").(int64)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	sessionID, err := strconv.ParseInt(req.SessionID, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid session_id")
		return
	}

	custID, err := strconv.ParseInt(req.CustomerID, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid customer_id")
		return
	}

	// Validasi customer ID dari body cocok dengan context
	if custID != customerID {
		response.Error(w, http.StatusForbidden, "forbidden")
		return
	}

	input := payment.ExecutePaymentInput{
		CustomerID: customerID,
		Total:      req.Total,
		SessionID:  sessionID,
	}

	result, err := h.paymentUC.Execute(r.Context(), input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := dto.PaymentInfoDto{
		Type:       result.Type,
		QrisString: result.QrisString,
		VaNumber:   result.VaNumber,
		BankName:   result.BankName,
		Amount:     result.Amount,
		ExpiredAt:  result.ExpiredAt,
	}

	response.Success(w, http.StatusOK, "Payment initiated", resp)
}
