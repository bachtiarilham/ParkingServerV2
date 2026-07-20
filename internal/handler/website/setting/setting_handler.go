package setting

import (
	"encoding/json"
	"net/http"

	"modulegue/core/response"
	dto "modulegue/internal/data/website/remote/dto/setting"
	mapper "modulegue/internal/data/website/remote/mapper/setting"
	addparlokuc "modulegue/internal/domain/web/usecase/setting"
)

type SettingHandler struct {
	addParlokUc         *addparlokuc.AddParlokUseCase
	registerJukirUc     *addparlokuc.RegisterJukirUseCase
	saveScheduleUc      *addparlokuc.SaveScheduleUseCase
	saveTarifUc         *addparlokuc.SaveTarifUseCase
	showSelectedJukirUc *addparlokuc.ShowSelectedJukirUseCase
	updateProfileUc     *addparlokuc.UpdateProfileUseCase
}

func NewSettingHandler(
	addParlokUc *addparlokuc.AddParlokUseCase,
	registerJukirUc *addparlokuc.RegisterJukirUseCase,
	saveScheduleUc *addparlokuc.SaveScheduleUseCase,
	saveTarifUc *addparlokuc.SaveTarifUseCase,
	showSelectedJukirUc *addparlokuc.ShowSelectedJukirUseCase,
	updateProfileUc *addparlokuc.UpdateProfileUseCase,
) *SettingHandler {
	return &SettingHandler{
		addParlokUc:         addParlokUc,
		registerJukirUc:     registerJukirUc,
		saveScheduleUc:      saveScheduleUc,
		saveTarifUc:         saveTarifUc,
		showSelectedJukirUc: showSelectedJukirUc,
		updateProfileUc:     updateProfileUc,
	}
}

func (h *SettingHandler) AddParlok(w http.ResponseWriter, r *http.Request) {
	var req dto.AddParlokRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	input := mapper.ToAddParlokRequestModel(&req)
	if input == nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	if err := h.addParlokUc.Execute(r.Context(), *input); err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal menambahkan parkir lokasi")
		return
	}

	response.Success(w, http.StatusOK, "parkir lokasi berhasil ditambahkan", nil)
}

func (h *SettingHandler) RegisterJukir(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	input := mapper.ToRegisterRequestModel(&req)
	if input == nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	if err := h.registerJukirUc.Execute(r.Context(), *input); err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal register jukir")
		return
	}

	response.Success(w, http.StatusOK, "jukir berhasil diregister", nil)
}

func (h *SettingHandler) SaveSchedule(w http.ResponseWriter, r *http.Request) {
	var req dto.SaveScheduleRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	input := mapper.ToSaveScheduleRequestModel(&req)
	if input == nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	if err := h.saveScheduleUc.Execute(r.Context(), *input); err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal menyimpan schedule jukir")
		return
	}

	response.Success(w, http.StatusOK, "schedule jukir berhasil disimpan", nil)
}

func (h *SettingHandler) SaveTarif(w http.ResponseWriter, r *http.Request) {
	var req dto.SaveTarifRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	input := mapper.ToSaveTarifRequestModel(&req)
	if input == nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	if err := h.saveTarifUc.Execute(r.Context(), *input); err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal menyimpan tarif")
		return
	}

	response.Success(w, http.StatusOK, "tarif berhasil disimpan", nil)
}

func (h *SettingHandler) ShowSelectedJukir(w http.ResponseWriter, r *http.Request) {
	var req dto.ShowSelectedJukirRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	input := mapper.ToShowSelectedJukirRequestModel(&req)
	if input == nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	resp, err := h.showSelectedJukirUc.Execute(r.Context(), *input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat jukir")
		return
	}

	response.Success(w, http.StatusOK, "data jukir berhasil dimuat", mapper.ToShowSelectedJukirResponseDto(resp))
}

func (h *SettingHandler) UpdateProfil(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateProfilRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	input := mapper.ToUpdateProfilRequestModel(&req)
	if input == nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	if err := h.updateProfileUc.Execute(r.Context(), *input); err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal update profil")
		return
	}

	response.Success(w, http.StatusOK, "profil berhasil diperbarui", nil)
}
