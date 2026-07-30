package paymentgate

import (
	"encoding/json"
	"net/http"

	"modulegue/core/response"
	"modulegue/internal/middleware"

	reqdto "modulegue/internal/data/mobile/remote/dto/payment_gate"
	mapper "modulegue/internal/data/mobile/remote/mapper/payment_gate"

	uc "modulegue/internal/domain/mobile/usecase/payment_gate"
)

type PaymentGateHandler struct {
	tfuc     *uc.PayTransferUseCase
	topupuc  *uc.PayTopUpUseCase
	memberuc *uc.PayMembershipUseCase
	parkiruc *uc.PayParkingUseCase
	cashuc   *uc.PayCashParkirUseCase
}

func NewPaymentGateHandler(
	tfuc *uc.PayTransferUseCase,
	topupuc *uc.PayTopUpUseCase,
	memberuc *uc.PayMembershipUseCase,
	parkiruc *uc.PayParkingUseCase,
	cashuc *uc.PayCashParkirUseCase,

) *PaymentGateHandler {
	return &PaymentGateHandler{
		tfuc:     tfuc,
		topupuc:  topupuc,
		memberuc: memberuc,
		parkiruc: parkiruc,
		cashuc:   cashuc,
	}
}

// Endpoint: POST /api/v2/linespot/checkout
func (h *PaymentGateHandler) Execute(w http.ResponseWriter, r *http.Request) {
	userID, okUserId := middleware.UserIDFromContext(r.Context())
	if !okUserId {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req reqdto.PayRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	input := mapper.ToPayRequestModel(&req)
	input.UserID = userID

	switch req.PaymentType {
	case "PARKIR":
		result, err := h.parkiruc.Execute(r.Context(), *input)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		response.Success(w, http.StatusOK, "Pembayaran diproses", mapper.ToPayResponseDto(result))
	case "TRANSFER":
		result, err := h.tfuc.Execute(r.Context(), *input)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		response.Success(w, http.StatusOK, "Transfer diproses", mapper.ToPayResponseDto(result))
	case "TOPUP":
		result, err := h.topupuc.Execute(r.Context(), *input)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		response.Success(w, http.StatusOK, "Topup diproses", mapper.ToPayResponseDto(result))
	case "MEMBERSHIP":
		result, err := h.memberuc.Execute(r.Context(), *input)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		response.Success(w, http.StatusOK, "Membership diproses", mapper.ToPayResponseDto(result))
	case "PARKIRCASH":
		result, err := h.cashuc.Execute(r.Context(), *input)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		response.Success(w, http.StatusOK, "Pembayaran Cash Berhasil", mapper.ToPayResponseDto(result))
	default:
		response.Error(w, http.StatusBadRequest, "type code error")
	}
}
