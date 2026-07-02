package subscription

type PromoModel struct {
	SNk   []string             `json:"snk,omitempty"`
	Promo []PromoTerpilihModel `json:"promo,omitempty"`
}
