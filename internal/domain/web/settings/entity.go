package settings

type TopFilters struct {
	Zones    []string `json:"zones"`
	Dates    string   `json:"dates"`
	Officers []string `json:"officers"`
}

type ExportQueueItem struct {
	Label  string `json:"label"`
	Status string `json:"status"`
	Note   string `json:"note"`
}

type AlertRuleItem struct {
	Title     string `json:"title"`
	Threshold string `json:"threshold"`
	Source    string `json:"source"`
	PIC       string `json:"pic"`
}

type ShiftTemplateItem struct {
	Label   string `json:"label"`
	Hours   string `json:"hours"`
	UseCase string `json:"useCase"`
}

type AdminRoleItem struct {
	Role   string `json:"role"`
	Access string `json:"access"`
	Owner  string `json:"owner"`
}

type NotificationItem struct {
	Channel  string `json:"channel"`
	Trigger  string `json:"trigger"`
	Response string `json:"response"`
}

type PaymentMethodItem struct {
	Label   string `json:"label"`
	Enabled bool   `json:"enabled"`
	Icon    string `json:"icon"`
}

type DefaultTariffItem struct {
	VehicleType string `json:"vehicleType"`
	FirstHour   int64  `json:"firstHour"`
	NextHour    int64  `json:"nextHour"`
	MaxRate     int64  `json:"maxRate"`
}
