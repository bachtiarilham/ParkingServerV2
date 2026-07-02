package home

type HomeModel struct {
	Profile         *ProfileModel         `json:"profile"`
	CustomerSummary *CustomerSummaryModel `json:"customerSummary"`
	JukirSummary    *JukirSummaryModel    `json:"jukirSummary"`
	Events          *EventsModel          `json:"events"`
	News            *NewsModel            `json:"news"`
	Warnings        *WarningsModel        `json:"warnings"`
}
