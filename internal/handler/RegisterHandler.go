package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"modulegue/core/errorstring"
	"modulegue/core/response"
	"modulegue/internal/data/mobile/remote/dto"
	"modulegue/internal/domain/mobile/usecase"
)

type RegisterHandler struct {
	registerUc *usecase.RegisterUseCase
}

func NewRegisterHandler(
	registerUc *usecase.RegisterUseCase,
) *RegisterHandler {
	return &RegisterHandler{
		registerUc: registerUc,
	}
}

func (h *RegisterHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	// --- Mapping: DTO -> UseCase Input ---
	input := usecase.RegisterInput{
		FullName: req.FullName,
		Nik:      req.NIK,
		Phone:    req.Phone,
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}

	// --- Panggil UseCase ---
	result, err := h.registerUc.Execute(r.Context(), input) // Asumsikan registerUC adalah *auth.RegisterUseCase
	if err != nil {
		// Log error jika perlu
		switch {
		case errors.Is(err, errorstring.ErrEmailAlreadyExists):
			response.Error(w, http.StatusConflict, "email sudah terdaftar")
		case errors.Is(err, errorstring.ErrUsernameAlreadyExists):
			response.Error(w, http.StatusConflict, "username sudah digunakan")
		default:
			// Error lainnya
			response.Error(w, http.StatusInternalServerError, "terjadi kesalahan internal")
		}
		return
	}

	// --- Mapping: UseCase Output -> DTO ---
	resp := dto.RegisterRe   sponseDto{
		Message: result.Message,
		User:  result.User,
	}

	response.Success(w, http.StatusCreated, "registrasi berhasil", resp)
}
