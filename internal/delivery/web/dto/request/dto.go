package request

import (
	"modulegue/internal/domain/web/location"
	"modulegue/internal/domain/web/settings"
)

type UpdateLocationSettingsRequest struct {
	TariffMotor     int64  `json:"tariffMotor"`
	TariffMobil     int64  `json:"tariffMobil"`
	DismissalReason string `json:"dismissalReason"`
}

type SaveShiftTemplatesRequest struct {
	ShiftTemplates []location.ParkingShiftTemplate `json:"shiftTemplates"`
}

type UpdateOfficerStatusRequest struct {
	Status string `json:"status"`
}

type ApplyOfficerMutationRequest struct {
	OfficerID        string `json:"officerId"`
	TargetLocationID string `json:"targetLocationId"`
	TargetShiftID    string `json:"targetShiftId"`
	Note             string `json:"note"`
}

type UpdateSettingsOverviewRequest struct {
	AlertRuleItems        []settings.AlertRuleItem     `json:"alertRuleItems"`
	DefaultShiftTemplates []settings.ShiftTemplateItem `json:"defaultShiftTemplates"`
	DefaultTariffItems    []settings.DefaultTariffItem `json:"defaultTariffItems"`
	NotificationItems     []settings.NotificationItem  `json:"notificationItems"`
	PaymentMethodItems    []settings.PaymentMethodItem `json:"paymentMethodItems"`
}

type CreateDisputeCaseRequest struct {
	ReferenceEntityType string `json:"referenceEntityType"`
	ReferenceEntityID   int64  `json:"referenceEntityId"`
	CaseType            string `json:"caseType"`
	AssignedToUserID    int64  `json:"assignedToUserId"`
	ChangeNote          string `json:"changeNote"`
}

type UpdateDisputeCaseStatusRequest struct {
	Status           string `json:"status"`
	AssignedToUserID int64  `json:"assignedToUserId"`
	ChangeNote       string `json:"changeNote"`
}

type CreateRefundTransactionRequest struct {
	ReferenceEntityType string `json:"referenceEntityType"`
	ReferenceEntityID   int64  `json:"referenceEntityId"`
	PaymentEventID      int64  `json:"paymentEventId"`
	WalletID            int64  `json:"walletId"`
	RefundAmount        int64  `json:"refundAmount"`
	RefundReason        string `json:"refundReason"`
}

type UpdateRefundStatusRequest struct {
	Status       string `json:"status"`
	StatusReason string `json:"statusReason"`
}

type CreateClosingBatchRequest struct {
	LocationID                 int64  `json:"locationId"`
	ClosingDate                string `json:"closingDate"`
	ActualClosingBalanceAmount int64  `json:"actualClosingBalanceAmount"`
	ChangeNote                 string `json:"changeNote"`
}

type UpdateClosingStatusRequest struct {
	Status     string `json:"status"`
	ChangeNote string `json:"changeNote"`
}
