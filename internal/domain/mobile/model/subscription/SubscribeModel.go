package subscription

type SubscribeModel struct {
	StatusCard  *StatusCardModel   `json:"status_card,omitempty"`
	PackageCard []PackageCardModel `json:"package_card,omitempty"`
	Promo       []PromoModel       `json:"promo,omitempty"`
}
