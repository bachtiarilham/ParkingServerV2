package settings

import (
	"io"
	"net/http"

	"modulegue/core/response"
	model "modulegue/internal/domain/mobile/model/settings"
	usecase "modulegue/internal/domain/mobile/usecase/settings"
	middleware "modulegue/internal/middleware"
)

type ChangeProfileHandler struct {
	changeProfileUc *usecase.ChangeProfileUseCase
}

func NewChangeProfileHandler(
	changeProfileUc *usecase.ChangeProfileUseCase,
) *ChangeProfileHandler {
	return &ChangeProfileHandler{
		changeProfileUc: changeProfileUc,
	}
}

func (h *ChangeProfileHandler) Execute(w http.ResponseWriter, r *http.Request) {
	userId, okUserId := middleware.UserIDFromContext(r.Context())
	if !okUserId {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse multipart form (10 MB maximum memory)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "Gagal memproses form multipart: "+err.Error())
		return
	}

	var input model.SettingsModel
	input.UserId = userId

	// Get text form values
	if emailVal := r.FormValue("email"); emailVal != "" {
		input.Email = &emailVal
	}
	if noTelpVal := r.FormValue("no_telp"); noTelpVal != "" {
		input.NoTelp = &noTelpVal
	}

	// Get file upload
	file, fileHeader, err := r.FormFile("foto_profil")
	if err == nil {
		defer file.Close()
		fileBytes, readErr := io.ReadAll(file)
		if readErr != nil {
			response.Error(w, http.StatusInternalServerError, "Gagal membaca file foto profil")
			return
		}
		input.FotoProfilBytes = fileBytes
		input.FotoProfilName = fileHeader.Filename
	}

	if err := h.changeProfileUc.Execute(r.Context(), input); err != nil {
		response.Error(w, http.StatusInternalServerError, "Gagal memperbarui profil: "+err.Error())
		return
	}

	response.Success(w, http.StatusOK, "Pembaruan profil berhasil", nil)
}
