package customer_membership

import "context"

type Repository interface {
	GetCustomerMembershipByUserID(ctx context.Context, userID int64) (*CustomerMembership, error)
	UpdateCustomerMembershipPoolBalance(ctx context.Context, membershipID int64, newBalance int64) error
}
