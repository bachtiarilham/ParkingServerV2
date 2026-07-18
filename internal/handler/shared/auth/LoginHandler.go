package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"modulegue/core/errorstring"
	"modulegue/core/response"
	dto "modulegue/internal/data/shared/remote/dto/auth"
	mapper "modulegue/internal/data/shared/remote/mapper/auth"
	usecase "modulegue/internal/domain/shared/usecase/auth"
	"modulegue/internal/middleware"
)

type LoginHandler struct {
	loginUc *usecase.LoginUseCase
}

func NewLoginHandler(
	loginUc *usecase.LoginUseCase,
) *LoginHandler {
	return &LoginHandler{
		loginUc: loginUc,
	}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	req.Identity = normalizeIdentity(req.Identity)
	req.Password = strings.TrimSpace(req.Password)

	requestModel := mapper.ToLoginReqModel(&req)
	if requestModel == nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}
	requestModel.DeviceName = middleware.UserAgent(r)
	requestModel.DeviceId = middleware.ClientIP(r)

	token, lol, err := h.loginUc.Execute(r.Context(), *requestModel)
	if err != nil {
		if errors.Is(err, errorstring.ErrInvalidCredentials) {
			response.Error(w, http.StatusUnauthorized, "identity atau password salah")
			return
		}
		if errors.Is(err, usecase.ErrInvalidInput) {
			response.Error(w, http.StatusBadRequest, "identity dan password wajib diisi")
			return
		}

		response.Error(w, http.StatusInternalServerError, "terjadi kesalahan internal")
		return
	}

	resp := mapper.ToLoginRespDto(token, lol.RoleId)
	if resp == nil {
		response.Error(w, http.StatusInternalServerError, "terjadi kesalahan internal")
		return
	}

	response.Success(w, http.StatusOK, "login berhasil", resp)
}

func normalizeIdentity(identity string) string {
	identity = strings.TrimSpace(identity)
	if strings.Contains(identity, "@") {
		return strings.ToLower(identity)
	}
	return identity
}
