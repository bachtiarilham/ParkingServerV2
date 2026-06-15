package handler

import (
	dto "modulegue/internal/delivery/mobile/customer/dto"
	middleware "modulegue/internal/middleware"
	home "modulegue/internal/usecase/home_customer"
	"modulegue/pkg/response"
	"net/http"
	"strconv"
)

type HomeHandler struct {
	getDashboardUC *home.GetDashboardUseCase
}

func NewHomeHandler(getDashboardUC *home.GetDashboardUseCase) *HomeHandler {
	return &HomeHandler{getDashboardUC: getDashboardUC}
}

// Endpoint: GET /api/v2/linespot/home?customer_id=123
func (h *HomeHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {

	userID, ok := middleware.UserIDFromContext(r.Context()) // <-- Gunakan fungsi helper
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	input := home.GetDashboardInput{
		UserID: userID,
		Limit:  5,
		Offset: 0,
	}
	result, err := h.getDashboardUC.Execute(r.Context(), input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat dashboard")
		return
	}

	// Mapping ke DTO
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
		}
	}
	for _, ev := range result.Events {
		resp.Events = append(resp.Events, dto.EventDto{
			Id:          strconv.FormatInt(ev.ID, 10),
			Title:       ev.Title,
			Description: ev.Description,
			Date:        ev.Date.Format("2006-01-02T15:04:05Z"),
			ImageUrl:    ev.ImageURL,
			Tag:         "EVENT",
		})
	}
	for _, nw := range result.News {
		resp.News = append(resp.News, dto.NewsDto{
			Id:          strconv.FormatInt(nw.ID, 10),
			Title:       nw.Title,
			Description: nw.Description,
			Date:        nw.Date.Format("2006-01-02T15:04:05Z"),
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

	response.Success(w, http.StatusOK, "Dashboard dimuat", resp)
}
