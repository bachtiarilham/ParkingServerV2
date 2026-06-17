package dto

import (
	"modulegue/internal/domain/web/settings"
)

type SettingsOverview struct {
	AlertRuleItems        []settings.AlertRuleItem     `json:"alertRuleItems"`
	DefaultShiftTemplates []settings.ShiftTemplateItem `json:"defaultShiftTemplates"`
	DefaultTariffItems    []settings.DefaultTariffItem `json:"defaultTariffItems"`
	AdminRoleItems        []settings.AdminRoleItem     `json:"adminRoleItems"`
	NotificationItems     []settings.NotificationItem  `json:"notificationItems"`
	PaymentMethodItems    []settings.PaymentMethodItem `json:"paymentMethodItems"`
}
