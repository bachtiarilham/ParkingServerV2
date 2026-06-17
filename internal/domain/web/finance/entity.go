package finance

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
