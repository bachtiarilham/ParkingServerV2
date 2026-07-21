package riwayat

type DetilParkirModel struct {
	Tanggal           string `json:"tanggal"`
	TransactionID     int64  `json:"transaction_id"`
	TransactionCode   string `json:"transaction_code"`
	SessionID         int64  `json:"session_id"`
	PlateNumber       string `json:"plate_number"`
	VehicleTypeID     int64  `json:"vehicle_type_id"`
	VehicleTypeCode   string `json:"vehicle_type_code"`
	VehicleTypeName   string `json:"vehicle_type_name"`
	PaymentMethodID   int64  `json:"payment_method_id"`
	PaymentMethodCode string `json:"payment_method_code"`
	PaymentMethodName string `json:"payment_method_name"`
	LocationID        int64  `json:"location_id"`
	LocationName      string `json:"location_name"`
	LocationAddress   string `json:"location_address"`
	AreaID            int64  `json:"area_id,omitempty"`
	AreaName          string `json:"area_name,omitempty"`
	ZoneID            int64  `json:"zone_id,omitempty"`
	ZoneName          string `json:"zone_name,omitempty"`
	BaseAmount        int64  `json:"base_amount"`
	DiscountAmount    int64  `json:"discount_amount"`
	FinalAmount       int64  `json:"final_amount"`
	CompanyShare      int64  `json:"company_share"`
	JukirShare        int64  `json:"jukir_share"`
	TaxAmount         int64  `json:"tax_amount"`
	FeeAmount         int64  `json:"fee_amount"`
	TransactionStatus string `json:"transaction_status"`
	OperationType     string `json:"operation_type"`
	OccurredAt        string `json:"occurred_at"`
	PaidAt            string `json:"paid_at"`
	CreatedAt         string `json:"created_at"`
}
