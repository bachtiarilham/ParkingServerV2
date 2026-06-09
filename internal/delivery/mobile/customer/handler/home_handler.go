package handler

import (
	// "encoding/json"
	dto "modulegue/internal/delivery/mobile/customer/dto"
	"modulegue/internal/usecase/home"
	"modulegue/pkg/response"
	"net/http"
	"strconv"
)

type HomeHandler struct {
	getDashboardUC *home.GetDashboardUseCase
}

func NewHomeHandler(getDashboardUC *home.GetDashboardUseCase) *HomeHandler {
	return &HomeHandler{
		getDashboardUC: getDashboardUC,
	}
}

func (h *HomeHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	// Ambil userID dari context (misalnya dari middleware JWT)
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Ambil query params untuk pagination (opsional)
	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")
	limit := 5  // Default
	offset := 0 // Default
	if limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	if offsetParam != "" {
		if parsedOffset, err := strconv.Atoi(offsetParam); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// --- Mapping: HTTP Request (params) -> UseCase Input ---
	input := home.GetDashboardInput{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	}

	// --- Panggil UseCase ---
	result, err := h.getDashboardUC.Execute(r.Context(), input)
	if err != nil {
		// Log error jika perlu
		response.Error(w, http.StatusInternalServerError, "Terjadi kesalahan internal")
		return
	}

	// --- Mapping: UseCase Output -> DTO Response ---
	resp := dto.HomeResponse{
		Profile:  nil,
		Summary:  nil,
		Events:   []dto.EventDto{},
		News:     []dto.NewsDto{},
		Warnings: nil,
	}

	if result.Profile != nil {
		resp.Profile = &dto.ProfileDto{
			Id:   strconv.FormatInt(result.Profile.ID, 10),
			Name: result.Profile.Name,
		}
	}

	if result.Summary != nil {
		resp.Summary = &dto.SummaryDto{
			Saldo: result.Summary.Saldo,
			// ExpiredDate: result.Summary.ExpiredDate.Format(time.RFC3339), // Jika ExpiredDate digunakan
		}
	}

	for _, ev := range result.Events {
		resp.Events = append(resp.Events, dto.EventDto{
			Id:          strconv.FormatInt(ev.ID, 10),
			Title:       ev.Title,
			Description: ev.Description,
			Date:        ev.Date.Format("2006-01-02T15:04:05Z07:00"), // Format ISO8601
			ImageUrl:    ev.ImageURL,
			Tag:         "EVENT",
		})
	}

	for _, nw := range result.News {
		resp.News = append(resp.News, dto.NewsDto{
			Id:          strconv.FormatInt(nw.ID, 10),
			Title:       nw.Title,
			Description: nw.Description,
			Date:        nw.Date.Format("2006-01-02T15:04:05Z07:00"), // Format ISO8601
			ImageUrl:    nw.ImageURL,
			Tag:         "NEWS",
		})
	}

	if result.Warnings != nil {
		resp.Warnings = &dto.WarningsDto{
			Profile: result.Warnings.Profile,
			Parking: result.Warnings.Parking,
			Finance: result.Warnings.Finance,
		}
	}

	response.Success(w, http.StatusOK, "Dashboard loaded", resp)
}
