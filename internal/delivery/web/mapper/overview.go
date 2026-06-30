package mapper

import (
	"modulegue/internal/delivery/web/dto"
	"modulegue/internal/domain/dashboard"
)

func FromDashboardOverview(src dashboard.DashboardOverview) dto.DashboardOverview {
	return dto.DashboardOverview{
		DashboardStats:        fromStatCards(src.DashboardStats),
		DashboardTransactions: fromRowItems(src.DashboardTransactions),
		FieldOfficers:         fromRowItems(src.FieldOfficers),
		DashboardAlerts:       fromAlertItems(src.DashboardAlerts),
		HourlyTraffic:         fromHourlyTraffic(src.HourlyTraffic),
		RevenueByLocation:     fromLocationMetrics(src.RevenueByLocation),
		OccupancyByLocation:   fromLocationMetrics(src.OccupancyByLocation),
		ComparisonMetrics:     fromComparisonMetrics(src.ComparisonMetrics),
		ParkingHeatmap:        fromHeatmapPoints(src.ParkingHeatmap),
		PriorityActions:       fromActionItems(src.PriorityActions),
		ParkingLocations:      fromParkingLocations(src.ParkingLocations),
	}
}

func FromMonitoringOverview(src dashboard.MonitoringOverview) dto.MonitoringOverview {
	return dto.MonitoringOverview{
		TopFilters:            dto.TopFilters{Zones: append([]string(nil), src.TopFilters.Zones...), Dates: src.TopFilters.Dates, Officers: append([]string(nil), src.TopFilters.Officers...)},
		MonitoringZones:       fromRowItems(src.MonitoringZones),
		ParkingLocations:      fromParkingLocations(src.ParkingLocations),
		ParkingOfficerOptions: fromParkingOfficerOptions(src.ParkingOfficerOptions),
	}
}

func FromOfficerOverview(src dashboard.OfficerOverview) dto.OfficerOverview {
	return dto.OfficerOverview{
		OfficerStats:          fromStatCards(src.OfficerStats),
		ParkingOfficerOptions: fromParkingOfficerOptions(src.ParkingOfficerOptions),
		ParkingLocations:      fromParkingLocations(src.ParkingLocations),
	}
}

func FromTransactionsOverview(src dashboard.TransactionsOverview) dto.TransactionsOverview {
	return dto.TransactionsOverview{
		TransactionStats:      fromStatCards(src.TransactionStats),
		TransactionRows:       fromRowItems(src.TransactionRows),
		PaymentBreakdownItems: fromPaymentBreakdownItems(src.PaymentBreakdownItems),
		TransactionIssueItems: fromTransactionIssueItems(src.TransactionIssueItems),
		ExportQueueItems:      fromExportQueueItems(src.ExportQueueItems),
	}
}

func FromSettingsOverview(src dashboard.SettingsOverview) dto.SettingsOverview {
	return dto.SettingsOverview{
		AlertRuleItems:        fromAlertRuleItems(src.AlertRuleItems),
		DefaultShiftTemplates: fromShiftTemplateItems(src.DefaultShiftTemplates),
		DefaultTariffItems:    fromDefaultTariffItems(src.DefaultTariffItems),
		AdminRoleItems:        fromAdminRoleItems(src.AdminRoleItems),
		NotificationItems:     fromNotificationItems(src.NotificationItems),
		PaymentMethodItems:    fromPaymentMethodItems(src.PaymentMethodItems),
	}
}

func FromAuthEnvelope(src dashboard.AuthEnvelope) dto.AuthEnvelope {
	return dto.AuthEnvelope{
		User: dto.AuthUser{
			UserID:     src.User.UserID,
			FullName:   src.User.FullName,
			Phone:      src.User.Phone,
			Email:      src.User.Email,
			Username:   src.User.Username,
			Role:       src.User.Role,
			IsVerified: src.User.IsVerified,
		},
		Tokens: dto.TokenSet{
			AccessToken:      src.Tokens.AccessToken,
			RefreshToken:     src.Tokens.RefreshToken,
			TokenType:        src.Tokens.TokenType,
			ExpiresInSeconds: src.Tokens.ExpiresInSeconds,
		},
	}
}

func fromStatCards(src []dashboard.StatCard) []dto.StatCard {
	out := make([]dto.StatCard, 0, len(src))
	for _, item := range src {
		out = append(out, dto.StatCard(item))
	}
	return out
}

func fromRowItems(src []dashboard.RowItem) []dto.RowItem {
	out := make([]dto.RowItem, 0, len(src))
	for _, item := range src {
		out = append(out, dto.RowItem(item))
	}
	return out
}

func fromAlertItems(src []dashboard.AlertItem) []dto.AlertItem {
	out := make([]dto.AlertItem, 0, len(src))
	for _, item := range src {
		out = append(out, dto.AlertItem(item))
	}
	return out
}

func fromHourlyTraffic(src []dashboard.HourlyTrafficPoint) []dto.HourlyTrafficPoint {
	out := make([]dto.HourlyTrafficPoint, 0, len(src))
	for _, item := range src {
		out = append(out, dto.HourlyTrafficPoint(item))
	}
	return out
}

func fromLocationMetrics(src []dashboard.LocationMetric) []dto.LocationMetric {
	out := make([]dto.LocationMetric, 0, len(src))
	for _, item := range src {
		out = append(out, dto.LocationMetric(item))
	}
	return out
}

func fromComparisonMetrics(src []dashboard.ComparisonMetric) []dto.ComparisonMetric {
	out := make([]dto.ComparisonMetric, 0, len(src))
	for _, item := range src {
		out = append(out, dto.ComparisonMetric(item))
	}
	return out
}

func fromHeatmapPoints(src []dashboard.HeatmapPoint) []dto.HeatmapPoint {
	out := make([]dto.HeatmapPoint, 0, len(src))
	for _, item := range src {
		out = append(out, dto.HeatmapPoint(item))
	}
	return out
}

func fromActionItems(src []dashboard.ActionItem) []dto.ActionItem {
	out := make([]dto.ActionItem, 0, len(src))
	for _, item := range src {
		out = append(out, dto.ActionItem(item))
	}
	return out
}

func fromParkingShiftTemplates(src []dashboard.ParkingShiftTemplate) []dto.ParkingShiftTemplate {
	out := make([]dto.ParkingShiftTemplate, 0, len(src))
	for _, item := range src {
		out = append(out, dto.ParkingShiftTemplate(item))
	}
	return out
}

func fromParkingLocations(src []dashboard.ParkingLocation) []dto.ParkingLocation {
	out := make([]dto.ParkingLocation, 0, len(src))
	for _, item := range src {
		out = append(out, dto.ParkingLocation{
			ID:                item.ID,
			Name:              item.Name,
			Zone:              item.Zone,
			Address:           item.Address,
			Lat:               item.Lat,
			Lng:               item.Lng,
			OfficerName:       item.OfficerName,
			OfficerShiftStart: item.OfficerShiftStart,
			OfficerShiftEnd:   item.OfficerShiftEnd,
			OfficerStatus:     item.OfficerStatus,
			DismissalReason:   item.DismissalReason,
			TariffMotor:       item.TariffMotor,
			TariffMobil:       item.TariffMobil,
			Motorcycles:       item.Motorcycles,
			Cars:              item.Cars,
			Officers:          item.Officers,
			OccupancyLabel:    item.OccupancyLabel,
			ShiftTemplates:    fromParkingShiftTemplates(item.ShiftTemplates),
		})
	}
	return out
}

func fromParkingOfficerOptions(src []dashboard.ParkingOfficerOption) []dto.ParkingOfficerOption {
	out := make([]dto.ParkingOfficerOption, 0, len(src))
	for _, item := range src {
		out = append(out, dto.ParkingOfficerOption(item))
	}
	return out
}

func fromPaymentBreakdownItems(src []dashboard.PaymentBreakdownItem) []dto.PaymentBreakdownItem {
	out := make([]dto.PaymentBreakdownItem, 0, len(src))
	for _, item := range src {
		out = append(out, dto.PaymentBreakdownItem(item))
	}
	return out
}

func fromTransactionIssueItems(src []dashboard.TransactionIssueItem) []dto.TransactionIssueItem {
	out := make([]dto.TransactionIssueItem, 0, len(src))
	for _, item := range src {
		out = append(out, dto.TransactionIssueItem(item))
	}
	return out
}

func fromExportQueueItems(src []dashboard.ExportQueueItem) []dto.ExportQueueItem {
	out := make([]dto.ExportQueueItem, 0, len(src))
	for _, item := range src {
		out = append(out, dto.ExportQueueItem(item))
	}
	return out
}

func fromAlertRuleItems(src []dashboard.AlertRuleItem) []dto.AlertRuleItem {
	out := make([]dto.AlertRuleItem, 0, len(src))
	for _, item := range src {
		out = append(out, dto.AlertRuleItem(item))
	}
	return out
}

func fromShiftTemplateItems(src []dashboard.ShiftTemplateItem) []dto.ShiftTemplateItem {
	out := make([]dto.ShiftTemplateItem, 0, len(src))
	for _, item := range src {
		out = append(out, dto.ShiftTemplateItem(item))
	}
	return out
}

func fromDefaultTariffItems(src []dashboard.DefaultTariffItem) []dto.DefaultTariffItem {
	out := make([]dto.DefaultTariffItem, 0, len(src))
	for _, item := range src {
		out = append(out, dto.DefaultTariffItem(item))
	}
	return out
}

func fromAdminRoleItems(src []dashboard.AdminRoleItem) []dto.AdminRoleItem {
	out := make([]dto.AdminRoleItem, 0, len(src))
	for _, item := range src {
		out = append(out, dto.AdminRoleItem(item))
	}
	return out
}

func fromNotificationItems(src []dashboard.NotificationItem) []dto.NotificationItem {
	out := make([]dto.NotificationItem, 0, len(src))
	for _, item := range src {
		out = append(out, dto.NotificationItem(item))
	}
	return out
}

func fromPaymentMethodItems(src []dashboard.PaymentMethodItem) []dto.PaymentMethodItem {
	out := make([]dto.PaymentMethodItem, 0, len(src))
	for _, item := range src {
		out = append(out, dto.PaymentMethodItem(item))
	}
	return out
}
