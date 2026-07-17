package topup

import (
	"encoding/json"
	"net/http"
	"strings"

	"modulegue/core/response"
	dto "modulegue/internal/data/mobile/remote/dto/topup"
	mapper "modulegue/internal/data/mobile/remote/mapper/topup"
	usecase "modulegue/internal/domain/mobile/usecase/topup"
)

type MidtransSignatureVerifier interface {
	VerifySignature(orderID, statusCode, grossAmount, signature string) bool
}

type MidtransNotificationHandler struct {
	callbackUc        *usecase.TopUpCallbackUseCase
	signatureVerifier MidtransSignatureVerifier
}

func NewMidtransNotificationHandler(callbackUc *usecase.TopUpCallbackUseCase, signatureVerifier MidtransSignatureVerifier) *MidtransNotificationHandler {
	return &MidtransNotificationHandler{
		callbackUc:        callbackUc,
		signatureVerifier: signatureVerifier,
	}
}

func (h *MidtransNotificationHandler) Execute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method tidak diizinkan")
		return
	}

	var req dto.QrisCallbackRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "payload callback tidak valid")
		return
	}

	req.OrderID = strings.TrimSpace(req.OrderID)
	req.StatusCode = strings.TrimSpace(req.StatusCode)
	req.GrossAmount = strings.TrimSpace(req.GrossAmount)
	req.SignatureKey = strings.TrimSpace(req.SignatureKey)
	req.TransactionStatus = strings.TrimSpace(req.TransactionStatus)

	if req.OrderID == "" {
		response.Error(w, http.StatusBadRequest, "order_id wajib diisi")
		return
	}

	if h.signatureVerifier != nil && req.SignatureKey != "" {
		if !h.signatureVerifier.VerifySignature(req.OrderID, req.StatusCode, req.GrossAmount, req.SignatureKey) {
			response.Error(w, http.StatusUnauthorized, "signature callback tidak valid")
			return
		}
	}

	input := mapper.ToQrisCallbackRequestModel(&req)
	if input == nil {
		response.Error(w, http.StatusBadRequest, "payload callback tidak valid")
		return
	}

	result, err := h.callbackUc.Execute(r.Context(), *input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memproses callback")
		return
	}

	statusCode := http.StatusOK
	if strings.EqualFold(result.PaymentStatusCode, "PENDING") {
		statusCode = http.StatusAccepted
	}

	response.Success(w, statusCode, "callback diproses", map[string]any{
		"order_id":            req.OrderID,
		"transaction_status":  req.TransactionStatus,
		"payment_status_code": result.PaymentStatusCode,
	})
}
