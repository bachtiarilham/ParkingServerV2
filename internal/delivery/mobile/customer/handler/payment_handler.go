package handler

// type PaymentHandler struct {
// 	executePaymentUC *payment.ExecutePaymentUseCase
// }

// func NewPaymentHandler(executePaymentUC *payment.ExecutePaymentUseCase) *PaymentHandler {
// 	return &PaymentHandler{executePaymentUC: executePaymentUC}
// }

// // Endpoint: POST /api/v2/linespot/pay
// func (h *PaymentHandler) ExecutePayment(w http.ResponseWriter, r *http.Request) {
// 	var req dto.ExecutePaymentRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		response.Error(w, http.StatusBadRequest, "request tidak valid")
// 		return
// 	}

// 	// Ambil customer ID dari JWT context
// 	userID, ok := r.Context().Value("userID").(int64)
// 	if !ok {
// 		response.Error(w, http.StatusUnauthorized, "Unauthorized")
// 		return
// 	}

// 	// Validasi: customer_id dari body harus sama dengan userID dari JWT
// 	custID, err := strconv.ParseInt(req.CustomerID, 10, 64)
// 	if err != nil || custID != userID {
// 		response.Error(w, http.StatusForbidden, "customer_id tidak valid")
// 		return
// 	}

// 	sessionID, err := strconv.ParseInt(req.SessionID, 10, 64)
// 	if err != nil {
// 		response.Error(w, http.StatusBadRequest, "session_id tidak valid")
// 		return
// 	}

// 	input := payment.ExecutePaymentInput{
// 		CustomerID: userID,
// 		Total:      req.Total,
// 		SessionID:  sessionID,
// 	}

// 	result, err := h.executePaymentUC.Execute(r.Context(), input)
// 	if err != nil {
// 		response.Error(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	resp := dto.PaymentInfoDto{
// 		Type:       result.Type,
// 		QrisString: result.QrisString,
// 		VaNumber:   result.VaNumber,
// 		BankName:   result.BankName,
// 		Amount:     result.Amount,
// 		ExpiredAt:  result.ExpiredAt,
// 	}

// 	response.Success(w, http.StatusOK, "Pembayaran diproses", resp)
// }
