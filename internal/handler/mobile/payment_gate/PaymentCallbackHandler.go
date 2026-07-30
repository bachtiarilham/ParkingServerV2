package paymentgate

import (
	"encoding/json"
	"net/http"
	"strings"

	"modulegue/core/response"
	dto "modulegue/internal/data/mobile/remote/dto/payment_gate"
	mapper "modulegue/internal/data/mobile/remote/mapper/payment_gate"
	uc "modulegue/internal/domain/mobile/usecase/payment_gate"
)

type PaymentCallbackHandler struct {
	callbackUC *uc.PaymentCallbackUseCase
	statusUC   *uc.GetPaymentStatusUseCase
}

func NewPaymentCallbackHandler(
	callbackUC *uc.PaymentCallbackUseCase,
	statusUC *uc.GetPaymentStatusUseCase,
) *PaymentCallbackHandler {
	return &PaymentCallbackHandler{
		callbackUC: callbackUC,
		statusUC:   statusUC,
	}
}

// Endpoint: POST /api/v2/payment/midtrans/callback
func (h *PaymentCallbackHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req dto.CallbackRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	input := mapper.ToCallbackRequestModel(&req)

	if err := h.callbackUC.Execute(r.Context(), *input); err != nil {
		// Log the error but return HTTP 200 OK anyway to prevent Midtrans from retrying unnecessarily if it's a validation error
		// Or return 400/500 if we want it retried. Usually we return 200 for processed webhooks.
		response.Error(w, http.StatusOK, err.Error())
		return
	}

	response.Success(w, http.StatusOK, "Callback processed successfully", map[string]string{"status": "ok"})
}

// Endpoint: GET /api/v2/payment/status/{transaction_code}
func (h *PaymentCallbackHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract transaction code from path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		response.Error(w, http.StatusBadRequest, "Invalid transaction code")
		return
	}
	txCode := parts[len(parts)-1]

	result, err := h.statusUC.Execute(r.Context(), txCode)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, "Transaction status retrieved", result)
}
