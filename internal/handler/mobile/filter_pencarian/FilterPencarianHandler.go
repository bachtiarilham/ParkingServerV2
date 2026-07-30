package riwayat

import (
	"encoding/json"
	"net/http"

	"modulegue/core/response"
	"modulegue/core/utils"
	reqdto "modulegue/internal/data/mobile/remote/dto/filter_pencarian"
	reqmodel "modulegue/internal/domain/mobile/model/filter_pencarian"

	laporanmapper "modulegue/internal/data/mobile/remote/mapper/laporan"
	riwayatmapper "modulegue/internal/data/mobile/remote/mapper/riwayat"

	uc "modulegue/internal/domain/mobile/usecase/filter_pencarian"
	middleware "modulegue/internal/middleware"
)

type FilterPencarianHandler struct {
	laporanuc *uc.GetLaporanUseCase
	trxuc     *uc.GetRiwayatTransaksiUseCase
	memberuc  *uc.GetRiwayatMembershipUseCase
	parkiruc  *uc.GetRiwayatParkirUseCase
}

func NewFilterPencarianHandler(
	getRiwayatParkirUc *uc.GetRiwayatParkirUseCase,
	getRiwayatTransaksiUc *uc.GetRiwayatTransaksiUseCase,
	getRiwayatMembershipUc *uc.GetRiwayatMembershipUseCase,
	getLaporanUc *uc.GetLaporanUseCase,

) *FilterPencarianHandler {
	return &FilterPencarianHandler{
		parkiruc:  getRiwayatParkirUc,
		trxuc:     getRiwayatTransaksiUc,
		memberuc:  getRiwayatMembershipUc,
		laporanuc: getLaporanUc,
	}
}

// Endpoint: POST /api/v2/linespot/riwayat
func (h *FilterPencarianHandler) Execute(w http.ResponseWriter, r *http.Request) {
	userID, okUserId := middleware.UserIDFromContext(r.Context())
	if !okUserId {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req reqdto.FilterPencarianDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	startDate, err := utils.ParseISODate(req.StartDate)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "startDate tidak valid")
		return
	}

	endDate, err := utils.ParseISODate(req.EndDate)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "endDate tidak valid")
		return
	}

	input := reqmodel.FilterPencarianModel{
		UserID:         userID,
		SearchTypeCode: req.SearchTypeCode,
		StartDate:      startDate,
		EndDate:        endDate,
	}

	switch req.SearchTypeCode {
	case "PARKIR":
		result, err := h.parkiruc.Execute(r.Context(), input)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "gagal memuat riwayat")
			return
		}
		response.Success(w, http.StatusOK, "Riwayat dimuat", riwayatmapper.ToRiwayatParkirDto(result))
	case "TRANSAKSI":
		result, err := h.trxuc.Execute(r.Context(), input)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "gagal memuat riwayat")
			return
		}
		response.Success(w, http.StatusOK, "Riwayat dimuat", riwayatmapper.ToRiwayatTransaksiDto(result))
	case "MEMBERSHIP":
		result, err := h.memberuc.Execute(r.Context(), input)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "gagal memuat riwayat")
			return
		}
		response.Success(w, http.StatusOK, "Riwayat dimuat", riwayatmapper.ToRiwayatMembershipDto(result))
	case "LAPORAN":
		result, err := h.laporanuc.Execute(r.Context(), input)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "gagal memuat riwayat")
			return
		}
		response.Success(w, http.StatusOK, "Riwayat dimuat", laporanmapper.ToLaporanDto(result))
	default:
		response.Error(w, http.StatusOK, "type code error")

	}
}
