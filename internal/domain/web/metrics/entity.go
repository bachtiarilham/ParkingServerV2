package metrics

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

type HourlyTrafficPoint struct {
	Label  string `json:"label"`
	Masuk  int64  `json:"masuk"`
	Keluar int64  `json:"keluar"`
}

type ComparisonMetric struct {
	Label     string `json:"label"`
	Today     int64  `json:"today"`
	Yesterday int64  `json:"yesterday"`
	Unit      string `json:"unit"`
}

type LocationMetric struct {
	Name      string `json:"name"`
	Value     int64  `json:"value"`
	Secondary string `json:"secondary"`
	Tone      string `json:"tone"` // Warna indikator (green, blue, orange, gold)
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

type AlertItem struct {
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	Priority string `json:"priority"`
	Action   string `json:"action"`
}
