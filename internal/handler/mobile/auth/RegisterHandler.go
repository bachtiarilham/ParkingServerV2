package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"modulegue/core/errorstring"
	"modulegue/core/response"
	dto "modulegue/internal/data/mobile/remote/dto/auth"
	mapper "modulegue/internal/data/mobile/remote/mapper/auth"
	usecase "modulegue/internal/domain/mobile/usecase/auth"
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

	req.FullName = strings.TrimSpace(req.FullName)
	req.NIK = strings.TrimSpace(req.NIK)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	input := mapper.ToRegisterRequestModel(&req)
	if input == nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	err := h.registerUc.Execute(r.Context(), *input)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidInput):
			response.Error(w, http.StatusBadRequest, "field wajib belum lengkap")
		case errors.Is(err, errorstring.ErrEmailAlreadyExists), errors.Is(err, usecase.ErrEmailAlreadyExists):
			response.Error(w, http.StatusConflict, "email sudah terdaftar")
		case errors.Is(err, errorstring.ErrUsernameAlreadyExists), errors.Is(err, usecase.ErrUsernameAlreadyExists):
			response.Error(w, http.StatusConflict, "username sudah digunakan")
		case errors.Is(err, usecase.ErrPhoneAlreadyExists):
			response.Error(w, http.StatusConflict, "nomor telepon sudah digunakan")
		default:
			response.Error(w, http.StatusInternalServerError, "terjadi kesalahan internal")
		}
		return
	}

	response.Success(w, http.StatusCreated, "registrasi berhasil", nil)
}
