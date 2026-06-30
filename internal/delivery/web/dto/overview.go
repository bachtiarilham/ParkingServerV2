package dto

type StatCard struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Note  string `json:"note,omitempty"`
	Icon  string `json:"icon"`
	Tone  string `json:"tone"`
}

type RowItem struct {
	TransactionID  string `json:"transactionId,omitempty"`
	PaymentEventID string `json:"paymentEventId,omitempty"`
	LocationID     string `json:"locationId,omitempty"`
	Primary        string `json:"primary"`
	Secondary      string `json:"secondary,omitempty"`
	Status         string `json:"status,omitempty"`
	StatusTone     string `json:"statusTone,omitempty"`
	ValueA         string `json:"valueA,omitempty"`
	ValueB         string `json:"valueB,omitempty"`
	Location       string `json:"location,omitempty"`
	Action         string `json:"action,omitempty"`
	Price          string `json:"price,omitempty"`
	Time           string `json:"time,omitempty"`
	Note           string `json:"note,omitempty"`
}

type AlertItem struct {
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	Priority string `json:"priority"`
	Action   string `json:"action"`
}

type HourlyTrafficPoint struct {
	Label  string `json:"label"`
	Masuk  int64  `json:"masuk"`
	Keluar int64  `json:"keluar"`
}

type LocationMetric struct {
	Name      string `json:"name"`
	Value     int64  `json:"value"`
	Secondary string `json:"secondary"`
	Tone      string `json:"tone"`
}

type ComparisonMetric struct {
	Label     string `json:"label"`
	Today     int64  `json:"today"`
	Yesterday int64  `json:"yesterday"`
	Unit      string `json:"unit"`
}

type HeatmapPoint struct {
	Day   string `json:"day"`
	Hour  string `json:"hour"`
	Value int64  `json:"value"`
}

type ActionItem struct {
	LocationID     string `json:"locationId"`
	Location       string `json:"location"`
	Issue          string `json:"issue"`
	Recommendation string `json:"recommendation"`
	Href           string `json:"href"`
	Tone           string `json:"tone"`
}

type ParkingShiftTemplate struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Start string `json:"start"`
	End   string `json:"end"`
}

type ParkingLocation struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	Zone              string                 `json:"zone"`
	Address           string                 `json:"address"`
	Lat               float64                `json:"lat"`
	Lng               float64                `json:"lng"`
	OfficerName       string                 `json:"officerName"`
	OfficerShiftStart string                 `json:"officerShiftStart"`
	OfficerShiftEnd   string                 `json:"officerShiftEnd"`
	OfficerStatus     string                 `json:"officerStatus"`
	DismissalReason   string                 `json:"dismissalReason"`
	TariffMotor       int64                  `json:"tariffMotor"`
	TariffMobil       int64                  `json:"tariffMobil"`
	Motorcycles       int64                  `json:"motorcycles"`
	Cars              int64                  `json:"cars"`
	Officers          int64                  `json:"officers"`
	OccupancyLabel    string                 `json:"occupancyLabel"`
	ShiftTemplates    []ParkingShiftTemplate `json:"shiftTemplates"`
}

type ParkingOfficerOption struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Role              string `json:"role"`
	HomeZone          string `json:"homeZone"`
	Availability      string `json:"availability"`
	AvailabilityNote  string `json:"availabilityNote"`
	CurrentAssignment string `json:"currentAssignment"`
	CurrentLocationID string `json:"currentLocationId"`
	CurrentShiftID    string `json:"currentShiftId"`
	Status            string `json:"status"`
	DefaultShiftStart string `json:"defaultShiftStart"`
	DefaultShiftEnd   string `json:"defaultShiftEnd"`
	DefaultStatus     string `json:"defaultStatus"`
}

type TopFilters struct {
	Zones    []string `json:"zones"`
	Dates    string   `json:"dates"`
	Officers []string `json:"officers"`
}

type PaymentBreakdownItem struct {
	Label  string `json:"label"`
	Amount string `json:"amount"`
	Share  string `json:"share"`
	Tone   string `json:"tone"`
}

type TransactionIssueItem struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Action string `json:"action"`
	Tone   string `json:"tone"`
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

type DashboardOverview struct {
	DashboardStats        []StatCard           `json:"dashboardStats"`
	DashboardTransactions []RowItem            `json:"dashboardTransactions"`
	FieldOfficers         []RowItem            `json:"fieldOfficers"`
	DashboardAlerts       []AlertItem          `json:"dashboardAlerts"`
	HourlyTraffic         []HourlyTrafficPoint `json:"hourlyTraffic"`
	RevenueByLocation     []LocationMetric     `json:"revenueByLocation"`
	OccupancyByLocation   []LocationMetric     `json:"occupancyByLocation"`
	ComparisonMetrics     []ComparisonMetric   `json:"comparisonMetrics"`
	ParkingHeatmap        []HeatmapPoint       `json:"parkingHeatmap"`
	PriorityActions       []ActionItem         `json:"priorityActions"`
	ParkingLocations      []ParkingLocation    `json:"parkingLocations"`
}

type MonitoringOverview struct {
	TopFilters            TopFilters             `json:"topFilters"`
	MonitoringZones       []RowItem              `json:"monitoringZones"`
	ParkingLocations      []ParkingLocation      `json:"parkingLocations"`
	ParkingOfficerOptions []ParkingOfficerOption `json:"parkingOfficerOptions"`
}

type OfficerOverview struct {
	OfficerStats          []StatCard             `json:"officerStats"`
	ParkingOfficerOptions []ParkingOfficerOption `json:"parkingOfficerOptions"`
	ParkingLocations      []ParkingLocation      `json:"parkingLocations"`
}

type TransactionsOverview struct {
	TransactionStats      []StatCard             `json:"transactionStats"`
	TransactionRows       []RowItem              `json:"transactionRows"`
	PaymentBreakdownItems []PaymentBreakdownItem `json:"paymentBreakdownItems"`
	TransactionIssueItems []TransactionIssueItem `json:"transactionIssueItems"`
	ExportQueueItems      []ExportQueueItem      `json:"exportQueueItems"`
}

type SettingsOverview struct {
	AlertRuleItems        []AlertRuleItem     `json:"alertRuleItems"`
	DefaultShiftTemplates []ShiftTemplateItem `json:"defaultShiftTemplates"`
	DefaultTariffItems    []DefaultTariffItem `json:"defaultTariffItems"`
	AdminRoleItems        []AdminRoleItem     `json:"adminRoleItems"`
	NotificationItems     []NotificationItem  `json:"notificationItems"`
	PaymentMethodItems    []PaymentMethodItem `json:"paymentMethodItems"`
}

type UpdateSettingsOverviewRequest struct {
	AlertRuleItems        []AlertRuleItem     `json:"alertRuleItems"`
	DefaultShiftTemplates []ShiftTemplateItem `json:"defaultShiftTemplates"`
	DefaultTariffItems    []DefaultTariffItem `json:"defaultTariffItems"`
	NotificationItems     []NotificationItem  `json:"notificationItems"`
	PaymentMethodItems    []PaymentMethodItem `json:"paymentMethodItems"`
}

type UpdateLocationSettingsRequest struct {
	TariffMotor     int64  `json:"tariffMotor"`
	TariffMobil     int64  `json:"tariffMobil"`
	DismissalReason string `json:"dismissalReason"`
}

type SaveShiftTemplatesRequest struct {
	ShiftTemplates []ParkingShiftTemplate `json:"shiftTemplates"`
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

type DisputeCaseSummary struct {
	ID                  int64  `json:"id"`
	DisputeCaseCode     string `json:"disputeCaseCode"`
	ReferenceEntityType string `json:"referenceEntityType"`
	ReferenceEntityID   int64  `json:"referenceEntityId"`
	CaseType            string `json:"caseType"`
	Status              string `json:"status"`
	OpenedByUserID      int64  `json:"openedByUserId"`
	AssignedToUserID    int64  `json:"assignedToUserId"`
	OpenedAt            string `json:"openedAt"`
	ResolvedAt          string `json:"resolvedAt,omitempty"`
	ResolutionNote      string `json:"resolutionNote,omitempty"`
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

type RefundTransactionSummary struct {
	ID                    int64  `json:"id"`
	RefundTransactionCode string `json:"refundTransactionCode"`
	ReferenceEntityType   string `json:"referenceEntityType"`
	ReferenceEntityID     int64  `json:"referenceEntityId"`
	PaymentEventID        int64  `json:"paymentEventId"`
	WalletID              int64  `json:"walletId"`
	RefundAmount          int64  `json:"refundAmount"`
	CurrencyCode          string `json:"currencyCode"`
	RefundReason          string `json:"refundReason"`
	Status                string `json:"status"`
	RequestedByUserID     int64  `json:"requestedByUserId"`
	ApprovedByUserID      int64  `json:"approvedByUserId"`
	RequestedAt           string `json:"requestedAt"`
	ApprovedAt            string `json:"approvedAt,omitempty"`
	ProcessedAt           string `json:"processedAt,omitempty"`
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

type ClosingBatchSummary struct {
	ID                           int64  `json:"id"`
	ClosingBatchCode             string `json:"closingBatchCode"`
	LocationID                   int64  `json:"locationId"`
	ClosingDate                  string `json:"closingDate"`
	OpeningBalanceAmount         int64  `json:"openingBalanceAmount"`
	CashSalesAmount              int64  `json:"cashSalesAmount"`
	CashlessSalesAmount          int64  `json:"cashlessSalesAmount"`
	TopupAmount                  int64  `json:"topupAmount"`
	RefundAmount                 int64  `json:"refundAmount"`
	AdjustmentAmount             int64  `json:"adjustmentAmount"`
	ExpectedClosingBalanceAmount int64  `json:"expectedClosingBalanceAmount"`
	ActualClosingBalanceAmount   int64  `json:"actualClosingBalanceAmount"`
	VarianceAmount               int64  `json:"varianceAmount"`
	Status                       string `json:"status"`
	SubmittedByUserID            int64  `json:"submittedByUserId"`
	ReviewedByUserID             int64  `json:"reviewedByUserId"`
	ApprovedByUserID             int64  `json:"approvedByUserId"`
	SubmittedAt                  string `json:"submittedAt,omitempty"`
	ReviewedAt                   string `json:"reviewedAt,omitempty"`
	ApprovedAt                   string `json:"approvedAt,omitempty"`
}
