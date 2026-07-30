package filterpencarian

type FilterPencarianDto struct {
	SearchTypeCode string `json:"searchTypeCode"`
	StartDate      string `json:"startDate"`
	EndDate        string `json:"endDate"`
}
