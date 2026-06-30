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
		Greeting: &dto.GreetingDto{
			Title:       result.Greeting.Title,
			Subtitle:    result.Greeting.Subtitle,
			AvatarLabel: result.Greeting.AvatarLabel,
		},
		BalanceCard: &dto.BalanceCardDto{
			Label:        result.BalanceCard.Label,
			Amount:       result.BalanceCard.Amount,
			PrimaryCta:   result.BalanceCard.PrimaryCta,
			SecondaryCta: result.BalanceCard.SecondaryCta,
		},
		PremiumCard: &dto.PremiumCardDto{
			Title:       result.PremiumCard.Title,
			Description: result.PremiumCard.Description,
			CtaLabel:    result.PremiumCard.CtaLabel,
			Badge:       result.PremiumCard.Badge,
		},
		Shortcuts:        []dto.ShortcutDto{},
		RecentActivities: []dto.ActivityDto{},
		Promotions:       []dto.PromotionDto{},
		Profile:          nil,
		Summary:          nil,
		Events:           []dto.EventDto{},
		News:             []dto.NewsDto{},
		Warnings:         nil,
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
	for _, sc := range result.Shortcuts {
		resp.Shortcuts = append(resp.Shortcuts, dto.ShortcutDto{
			Title:    sc.Title,
			Icon:     sc.Icon,
			DeepLink: sc.DeepLink,
		})
	}
	for _, act := range result.RecentActivities {
		resp.RecentActivities = append(resp.RecentActivities, dto.ActivityDto{
			Title:       act.Title,
			Subtitle:    act.Subtitle,
			Status:      act.Status,
			ActionLabel: act.ActionLabel,
		})
	}
	for _, promo := range result.Promotions {
		resp.Promotions = append(resp.Promotions, dto.PromotionDto{
			Title:       promo.Title,
			Description: promo.Description,
			Badge:       promo.Badge,
		})
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
