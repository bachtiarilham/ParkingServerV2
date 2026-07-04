package payment

type PostPaymentParkingResponseDto struct {
	Title               *string                           `json:"title,omitempty"`
	SuccessTitle        *string                           `json:"success_title,omitempty"`
	SuccessDescription  *string                           `json:"success_description,omitempty"`
	TotalAmount         *string                           `json:"total_amount,omitempty"`
	PaymentStatus       *string                           `json:"payment_status,omitempty"`
	ReferenceNumber     *string                           `json:"reference_number,omitempty"`
	VerificationMessage *string                           `json:"verification_message,omitempty"`
	ThankYouTitle       *string                           `json:"thank_you_title,omitempty"`
	ThankYouDescription *string                           `json:"thank_you_description,omitempty"`
	DownloadLabel       *string                           `json:"download_label,omitempty"`
	BackToHomeLabel     *string                           `json:"back_to_home_label,omitempty"`
	Details             []PostPaymentParkingDetailItemDto `json:"details,omitempty"`
}

type PostPaymentParkingDetailItemDto struct {
	Label *string `json:"label,omitempty"`
	Value *string `json:"value,omitempty"`
}
