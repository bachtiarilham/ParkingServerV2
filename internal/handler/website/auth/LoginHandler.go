package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"modulegue/core/errorstring"
	"modulegue/core/response"
	dto "modulegue/internal/data/shared/remote/dto/auth"
	authmapper "modulegue/internal/data/shared/remote/mapper/auth"
	homemapper "modulegue/internal/data/website/remote/mapper/home"
	usecase "modulegue/internal/domain/web/usecase/login"
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

	requestModel := authmapper.ToLoginReqModel(&req)
	if requestModel == nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}
	requestModel.DeviceName = middleware.UserAgent(r)
	requestModel.DeviceId = middleware.ClientIP(r)

	_, result, err := h.loginUc.Execute(r.Context(), *requestModel)
	if err != nil {
		if errors.Is(err, errorstring.ErrInvalidCredentials) {
			response.Error(w, http.StatusUnauthorized, "identity atau password salah")
			return
		}
		if errors.Is(err, errors.New("tolong isi field yang wajib")) {
			response.Error(w, http.StatusBadRequest, "identity dan password wajib diisi")
			return
		}

		response.Error(w, http.StatusInternalServerError, "terjadi kesalahan internal")
		return
	}

	resp := homemapper.ToHomeResponseDto(result)
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
