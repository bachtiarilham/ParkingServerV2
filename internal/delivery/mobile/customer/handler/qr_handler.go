package handler

import (
	"modulegue/internal/delivery/mobile/customer/dto"
	"modulegue/internal/middleware"
	"modulegue/internal/usecase/qr"
	"modulegue/pkg/response"
	"net/http"
)

type QRHandler struct {
	generateQRUC *qr.GenerateQRUseCase
}

func NewQRHandler(generateQRUC *qr.GenerateQRUseCase) *QRHandler {
	return &QRHandler{
		generateQRUC: generateQRUC,
	}
}

// Endpoint: GET /api/v2/linespot/qr/generate
func (h *QRHandler) GenerateQR(w http.ResponseWriter, r *http.Request) {
	// Ambil user ID dari JWT context
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Input untuk usecase
	input := qr.GenerateQRInput{
		UserID: userID,
	}

	result, err := h.generateQRUC.Execute(r.Context(), input)
	if err != nil {
		// Log error jika perlu
		response.Error(w, http.StatusInternalServerError, "gagal membuat qr code")
		return
	}

	// Mapping ke DTO
	resp := dto.GenerateQRResponse{
		QRString: result.QRString,
	}

	response.Success(w, http.StatusOK, "QR Code dibuat", resp)
}
