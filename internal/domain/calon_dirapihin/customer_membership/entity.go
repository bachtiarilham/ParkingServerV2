package customer_membership

import "time"

type CustomerMembership struct {
	ID             int64     `json:"id"`
	CustomerUserID int64     `json:"customer_user_id"`
	PlanID         int       `json:"plan_id"`
	PoolBalance    int64     `json:"pool_balance"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}
