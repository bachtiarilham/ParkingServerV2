package dto

type HomeResponse struct {
	Greeting         *GreetingDto    `json:"greeting"`
	BalanceCard      *BalanceCardDto `json:"balanceCard"`
	PremiumCard      *PremiumCardDto `json:"premiumCard"`
	Shortcuts        []ShortcutDto   `json:"shortcuts"`
	RecentActivities []ActivityDto   `json:"recentActivities"`
	Promotions       []PromotionDto  `json:"promotions"`
	Profile          *ProfileDto     `json:"profile"`
	Summary          *SummaryDto     `json:"summary"`
	Events           []EventDto      `json:"events"`
	News             []NewsDto       `json:"news"`
	Warnings         *WarningsDto    `json:"warnings"`
}

type GreetingDto struct {
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	AvatarLabel string `json:"avatarLabel"`
}

type BalanceCardDto struct {
	Label        string `json:"label"`
	Amount       string `json:"amount"`
	PrimaryCta   string `json:"primaryCta"`
	SecondaryCta string `json:"secondaryCta"`
}

type PremiumCardDto struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CtaLabel    string `json:"ctaLabel"`
	Badge       string `json:"badge"`
}

type ShortcutDto struct {
	Title    string `json:"title"`
	Icon     string `json:"icon"`
	DeepLink string `json:"deeplink"`
}

type ActivityDto struct {
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Status      string `json:"status"`
	ActionLabel string `json:"actionLabel"`
}

type PromotionDto struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Badge       string `json:"badge"`
}

type ProfileDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type SummaryDto struct {
	Saldo       int64  `json:"saldo"`
	ExpiredDate string `json:"expiredDate"` // String untuk format ISO8601
}

type EventDto struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"` // String untuk format ISO8601
	ImageUrl    string `json:"imageurl"`
	Tag         string `json:"tag"`
}

type NewsDto struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"` // String untuk format ISO8601
	ImageUrl    string `json:"imageurl"`
	Tag         string `json:"tag"`
}

type WarningsDto struct {
	Profile string `json:"profile"`
	Parking string `json:"parking"`
	Finance string `json:"finance"`
}
