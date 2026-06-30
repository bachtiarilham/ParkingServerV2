package handler

import (
	"encoding/json"
	dto "modulegue/internal/delivery/mobile/jukir/dto"
	middleware "modulegue/internal/middleware"
	"modulegue/internal/usecase/payment"
	"modulegue/pkg/response"
	"net/http"
)

// Konstanta untuk role ID
const (
	ROLE_CUSTOMER_ID = 1 // Sesuaikan dengan ID di tabel system_role
	ROLE_JUKIR_ID    = 2 // Sesuaikan dengan ID di tabel system_role
)

type PaymentHandler struct {
	initiatePaymentUC *payment.InitiatePaymentUseCase
}

func NewPaymentHandler(initiatePaymentUC *payment.InitiatePaymentUseCase) *PaymentHandler {
	return &PaymentHandler{
		initiatePaymentUC: initiatePaymentUC,
	}
}

// Endpoint: POST /api/v2/linespot/pay (untuk initiate payment - digunakan oleh jukir app)
func (h *PaymentHandler) InitiatePayment(w http.ResponseWriter, r *http.Request) {
	var req dto.ExecutePaymentRequest // Atau gunakan struct DTO khusus untuk initiate payment jika berbeda
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	// Ambil JUKIR ID dari JWT context (ini adalah user yang sedang login - HARUS JUKIR)
	jukirID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// *** PENGECEKAN ROLE DI HANDLER MENGGUNAKAN HELPER ***
	roleID, ok := middleware.RoleIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized (role missing from token)")
		return
	}

	// Cek apakah role user adalah Jukir
	if roleID != ROLE_JUKIR_ID {
		response.Error(w, http.StatusForbidden, "akses ditolak: hanya jukir yang bisa mengakses endpoint ini")
		return
	}

	// Karena role valid, userID yang diambil sebelumnya (jukirID) adalah benar milik jukir yang sah

	// Validasi input sederhana
	if req.TransactionID == "" || req.PaymentMethod == "" {
		response.Error(w, http.StatusBadRequest, "transaction_id dan payment_method diperlukan")
		return
	}

	// Validasi customer_id dari body (harus angka dan positif)
	// customerID, err := strconv.ParseInt(req.CustomerID, 10, 64)
	// if err != nil {
	// 	response.Error(w, http.StatusBadRequest, "customer_id tidak valid")
	// 	return
	// }
	// Kita asumsikan req.CustomerID adalah string seperti di DTO Android
	// Kita gunakan req.CustomerID sebagai string untuk pencarian transaksi di usecase
	// UseCase nanti bisa mengonversinya ke int64 jika perlu atau mencari berdasarkan code

	// Input untuk usecase
	// Pastikan InitiatePaymentInput di usecase sesuai
	input := payment.InitiatePaymentInput{
		TransactionCode: req.TransactionID, // Kirim code transaksi (atau ID jika req.TransactionID adalah ID)
		CustomerID:      req.CustomerID,    // Customer yang transaksinya sedang diproses
		JukirID:         jukirID,           // Jukir yang melayani (didapat dari JWT dan validasi role)
		PaymentMethod:   req.PaymentMethod,
	}

	result, err := h.initiatePaymentUC.Execute(r.Context(), input)
	if err != nil {
		// Log error jika perlu
		switch {
		case err.Error() == "transaksi tidak ditemukan":
			response.Error(w, http.StatusNotFound, "transaksi tidak ditemukan")
		case err.Error() == "transaksi tidak dalam status pending/unpaid":
			response.Error(w, http.StatusBadRequest, "transaksi tidak bisa diproses pembayaran")
		case err == payment.ErrInsufficientJukirBalance:
			response.Error(w, http.StatusBadRequest, "saldo jukir tidak mencukupi untuk transaksi cash")
		case err == payment.ErrCustomerNotMember:
			response.Error(w, http.StatusBadRequest, "customer tidak memiliki membership aktif")
		case err == payment.ErrPaymentMethodUnsupported:
			response.Error(w, http.StatusBadRequest, "metode pembayaran tidak didukung")
		default:
			response.Error(w, http.StatusInternalServerError, "gagal menginisiasi pembayaran")
		}
		return
	}

	// Mapping ke DTO yang sesuai untuk response
	// Misalnya, gunakan dto.InitiatePaymentResponse dari file DTO kamu
	resp := dto.InitiatePaymentResponse{ // Ganti dengan struct DTO yang benar
		Type:       result.Type,
		QrisString: result.QrisString,
		VaNumber:   result.VaNumber,
		BankName:   result.BankName,
		Amount:     result.Amount,
		ExpiredAt:  result.ExpiredAt,
		Message:    result.Message,
	}

	response.Success(w, http.StatusOK, "Pembayaran diproses", resp)
}
