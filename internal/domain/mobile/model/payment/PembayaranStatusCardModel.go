package payment

type PembayaranStatusCardModel struct {
	Title     *string `json:"title,omitempty"`
	Message   *string `json:"message,omitempty"`
	IsSuccess *bool   `json:"is_success,omitempty"`
}
