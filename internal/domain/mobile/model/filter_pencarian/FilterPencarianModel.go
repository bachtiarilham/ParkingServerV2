package filterpencarian

import "time"

type FilterPencarianModel struct {
	UserID         int64     `json:"user_id"`
	SearchTypeCode string    `json:"searchTypeCode"`
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
}
