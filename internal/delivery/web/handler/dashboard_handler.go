package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"modulegue/internal/delivery/web/dto"
	"modulegue/internal/delivery/web/mapper"
	"modulegue/internal/domain/dashboard"
	"modulegue/internal/middleware"
	"modulegue/internal/usecase/web"
	"modulegue/pkg/response"
)

type DashboardHandler struct {
	dashboardUC *web.DashboardUseCase
}

func NewDashboardHandler(dashboardUC *web.DashboardUseCase) *DashboardHandler {
	return &DashboardHandler{dashboardUC: dashboardUC}
}

func (h *DashboardHandler) GetDashboardOverview(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := h.dashboardUC.GetDashboardOverview(r.Context(), userID)
	if err != nil {
		h.writeError(w, err, "gagal memuat dashboard")
		return
	}
	response.Success(w, http.StatusOK, "dashboard overview berhasil dimuat", mapper.FromDashboardOverview(result))
}

func (h *DashboardHandler) GetMonitoringOverview(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := h.dashboardUC.GetMonitoringOverview(r.Context(), userID)
	if err != nil {
		h.writeError(w, err, "gagal memuat monitoring")
		return
	}
	response.Success(w, http.StatusOK, "monitoring overview berhasil dimuat", mapper.FromMonitoringOverview(result))
}

func (h *DashboardHandler) GetOfficerOverview(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := h.dashboardUC.GetOfficerOverview(r.Context(), userID)
	if err != nil {
		h.writeError(w, err, "gagal memuat officer overview")
		return
	}
	response.Success(w, http.StatusOK, "officer overview berhasil dimuat", mapper.FromOfficerOverview(result))
}

func (h *DashboardHandler) GetTransactionsOverview(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := h.dashboardUC.GetTransactionsOverview(r.Context(), userID)
	if err != nil {
		h.writeError(w, err, "gagal memuat transactions overview")
		return
	}
	response.Success(w, http.StatusOK, "transactions overview berhasil dimuat", mapper.FromTransactionsOverview(result))
}

func (h *DashboardHandler) GetSettingsOverview(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := h.dashboardUC.GetSettingsOverview(r.Context(), userID)
	if err != nil {
		h.writeError(w, err, "gagal memuat settings overview")
		return
	}
	response.Success(w, http.StatusOK, "settings overview berhasil dimuat", mapper.FromSettingsOverview(result))
}

func (h *DashboardHandler) UpdateSettingsOverview(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req dto.UpdateSettingsOverviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}
	_, err := h.dashboardUC.UpdateSettingsOverview(r.Context(), userID, dashboard.UpdateSettingsOverviewRequest{
		AlertRuleItems:        toDashboardAlertRuleItems(req.AlertRuleItems),
		DefaultShiftTemplates: toDashboardShiftTemplateItems(req.DefaultShiftTemplates),
		DefaultTariffItems:    toDashboardDefaultTariffItems(req.DefaultTariffItems),
		NotificationItems:     toDashboardNotificationItems(req.NotificationItems),
		PaymentMethodItems:    toDashboardPaymentMethodItems(req.PaymentMethodItems),
	})
	if err != nil {
		h.writeError(w, err, "gagal memperbarui settings overview")
		return
	}
	response.Success(w, http.StatusOK, "settings overview berhasil diperbarui", nil)
}

func (h *DashboardHandler) UpdateLocationSettings(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req dto.UpdateLocationSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}
	_, err := h.dashboardUC.UpdateLocationSettings(r.Context(), userID, r.PathValue("id"), dashboard.UpdateLocationSettingsRequest{
		TariffMotor:     req.TariffMotor,
		TariffMobil:     req.TariffMobil,
		DismissalReason: req.DismissalReason,
	})
	if err != nil {
		h.writeError(w, err, "gagal memperbarui lokasi")
		return
	}
	response.Success(w, http.StatusOK, "lokasi berhasil diperbarui", nil)
}

func (h *DashboardHandler) SaveLocationShiftTemplates(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req dto.SaveShiftTemplatesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}
	templates := make([]dashboard.ParkingShiftTemplate, 0, len(req.ShiftTemplates))
	for _, item := range req.ShiftTemplates {
		templates = append(templates, dashboard.ParkingShiftTemplate(item))
	}
	_, err := h.dashboardUC.SaveLocationShiftTemplates(r.Context(), userID, r.PathValue("id"), templates)
	if err != nil {
		h.writeError(w, err, "gagal menyimpan shift template")
		return
	}
	response.Success(w, http.StatusOK, "shift template berhasil disimpan", nil)
}

func (h *DashboardHandler) UpdateOfficerStatus(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req dto.UpdateOfficerStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}
	_, err := h.dashboardUC.UpdateOfficerStatus(r.Context(), userID, r.PathValue("id"), req.Status)
	if err != nil {
		h.writeError(w, err, "gagal memperbarui status officer")
		return
	}
	response.Success(w, http.StatusOK, "status officer berhasil diperbarui", nil)
}

func (h *DashboardHandler) ApplyOfficerMutation(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req dto.ApplyOfficerMutationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}
	_, err := h.dashboardUC.ApplyOfficerMutation(r.Context(), userID, dashboard.ApplyOfficerMutationRequest{
		OfficerID:        req.OfficerID,
		TargetLocationID: req.TargetLocationID,
		TargetShiftID:    req.TargetShiftID,
		Note:             req.Note,
	})
	if err != nil {
		h.writeError(w, err, "gagal memproses mutasi officer")
		return
	}
	response.Success(w, http.StatusOK, "mutasi officer berhasil diproses", nil)
}

func (h *DashboardHandler) CreateDisputeCase(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req dto.CreateDisputeCaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}
	_, err := h.dashboardUC.CreateDisputeCase(r.Context(), userID, dashboard.CreateDisputeCaseRequest{
		ReferenceEntityType: req.ReferenceEntityType,
		ReferenceEntityID:   req.ReferenceEntityID,
		CaseType:            req.CaseType,
		AssignedToUserID:    req.AssignedToUserID,
		ChangeNote:          req.ChangeNote,
	})
	if err != nil {
		h.writeError(w, err, "gagal membuat dispute case")
		return
	}
	response.Success(w, http.StatusCreated, "dispute case berhasil dibuat", nil)
}

func (h *DashboardHandler) UpdateDisputeCaseStatus(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req dto.UpdateDisputeCaseStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}
	_, err := h.dashboardUC.UpdateDisputeCaseStatus(r.Context(), userID, r.PathValue("id"), dashboard.UpdateDisputeCaseStatusRequest{
		Status:           req.Status,
		AssignedToUserID: req.AssignedToUserID,
		ChangeNote:       req.ChangeNote,
	})
	if err != nil {
		h.writeError(w, err, "gagal memperbarui dispute case")
		return
	}
	response.Success(w, http.StatusOK, "dispute case berhasil diperbarui", nil)
}

func (h *DashboardHandler) CreateRefundTransaction(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req dto.CreateRefundTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}
	_, err := h.dashboardUC.CreateRefundTransaction(r.Context(), userID, dashboard.CreateRefundTransactionRequest{
		ReferenceEntityType: req.ReferenceEntityType,
		ReferenceEntityID:   req.ReferenceEntityID,
		PaymentEventID:      req.PaymentEventID,
		WalletID:            req.WalletID,
		RefundAmount:        req.RefundAmount,
		RefundReason:        req.RefundReason,
	})
	if err != nil {
		h.writeError(w, err, "gagal membuat refund transaction")
		return
	}
	response.Success(w, http.StatusCreated, "refund transaction berhasil dibuat", nil)
}

func (h *DashboardHandler) UpdateRefundTransactionStatus(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req dto.UpdateRefundStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}
	_, err := h.dashboardUC.UpdateRefundTransactionStatus(r.Context(), userID, r.PathValue("id"), dashboard.UpdateRefundStatusRequest{
		Status:       req.Status,
		StatusReason: req.StatusReason,
	})
	if err != nil {
		h.writeError(w, err, "gagal memperbarui refund transaction")
		return
	}
	response.Success(w, http.StatusOK, "refund transaction berhasil diperbarui", nil)
}

func (h *DashboardHandler) CreateClosingBatch(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req dto.CreateClosingBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}
	_, err := h.dashboardUC.CreateClosingBatch(r.Context(), userID, dashboard.CreateClosingBatchRequest{
		LocationID:                 req.LocationID,
		ClosingDate:                req.ClosingDate,
		ActualClosingBalanceAmount: req.ActualClosingBalanceAmount,
		ChangeNote:                 req.ChangeNote,
	})
	if err != nil {
		h.writeError(w, err, "gagal membuat closing batch")
		return
	}
	response.Success(w, http.StatusCreated, "closing batch berhasil dibuat", nil)
}

func (h *DashboardHandler) UpdateClosingBatchStatus(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req dto.UpdateClosingStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}
	_, err := h.dashboardUC.UpdateClosingBatchStatus(r.Context(), userID, r.PathValue("id"), dashboard.UpdateClosingStatusRequest{
		Status:     req.Status,
		ChangeNote: req.ChangeNote,
	})
	if err != nil {
		h.writeError(w, err, "gagal memperbarui closing batch")
		return
	}
	response.Success(w, http.StatusOK, "closing batch berhasil diperbarui", nil)
}

func (h *DashboardHandler) writeError(w http.ResponseWriter, err error, fallback string) {
	switch {
	case errors.Is(err, web.ErrForbiddenRole):
		response.Error(w, http.StatusForbidden, "role tidak diizinkan")
	default:
		response.Error(w, http.StatusInternalServerError, fallback)
	}
}

func toDashboardAlertRuleItems(src []dto.AlertRuleItem) []dashboard.AlertRuleItem {
	out := make([]dashboard.AlertRuleItem, 0, len(src))
	for _, item := range src {
		out = append(out, dashboard.AlertRuleItem{
			Title:     item.Title,
			Threshold: item.Threshold,
			Source:    item.Source,
			PIC:       item.PIC,
		})
	}
	return out
}

func toDashboardShiftTemplateItems(src []dto.ShiftTemplateItem) []dashboard.ShiftTemplateItem {
	out := make([]dashboard.ShiftTemplateItem, 0, len(src))
	for _, item := range src {
		out = append(out, dashboard.ShiftTemplateItem{
			Label:   item.Label,
			Hours:   item.Hours,
			UseCase: item.UseCase,
		})
	}
	return out
}

func toDashboardDefaultTariffItems(src []dto.DefaultTariffItem) []dashboard.DefaultTariffItem {
	out := make([]dashboard.DefaultTariffItem, 0, len(src))
	for _, item := range src {
		out = append(out, dashboard.DefaultTariffItem{
			VehicleType: item.VehicleType,
			FirstHour:   item.FirstHour,
			NextHour:    item.NextHour,
			MaxRate:     item.MaxRate,
		})
	}
	return out
}

func toDashboardNotificationItems(src []dto.NotificationItem) []dashboard.NotificationItem {
	out := make([]dashboard.NotificationItem, 0, len(src))
	for _, item := range src {
		out = append(out, dashboard.NotificationItem{
			Channel:  item.Channel,
			Trigger:  item.Trigger,
			Response: item.Response,
		})
	}
	return out
}

func toDashboardPaymentMethodItems(src []dto.PaymentMethodItem) []dashboard.PaymentMethodItem {
	out := make([]dashboard.PaymentMethodItem, 0, len(src))
	for _, item := range src {
		out = append(out, dashboard.PaymentMethodItem{
			Label:   item.Label,
			Enabled: item.Enabled,
			Icon:    item.Icon,
		})
	}
	return out
}
