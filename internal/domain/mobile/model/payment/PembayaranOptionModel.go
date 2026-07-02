package payment

type PembayaranOptionModel struct {
	Type     *string `json:"type,omitempty"`
	Title    *string `json:"title,omitempty"`
	Subtitle *string `json:"subtitle,omitempty"`
}
