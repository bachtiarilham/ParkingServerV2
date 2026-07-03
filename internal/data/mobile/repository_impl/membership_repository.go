package repository

import (
	"context"
	"database/sql"
	"fmt"
	"modulegue/internal/domain/customer_membership"
)

type MembershipRepository struct {
	db *sql.DB
}

func NewMembershipRepository(db *sql.DB) customer_membership.Repository {
	return &MembershipRepository{db: db}
}

func (r *MembershipRepository) GetCustomerMembershipByUserID(ctx context.Context, userID int64) (*customer_membership.CustomerMembership, error) {
	query := `
		SELECT id, customer_user_id, plan_id, pool_balance, start_date, end_date, status, created_at
		FROM customer_memberships
		WHERE customer_user_id = ? AND status = 'active' AND end_date > NOW()
		LIMIT 1
	`
	var m customer_membership.CustomerMembership
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&m.ID, &m.CustomerUserID, &m.PlanID, &m.PoolBalance, &m.StartDate, &m.EndDate, &m.Status, &m.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("active customer membership not found for user id: %d", userID)
		}
		return nil, fmt.Errorf("get customer membership by user id: %w", err)
	}
	return &m, nil
}

func (r *MembershipRepository) UpdateCustomerMembershipPoolBalance(ctx context.Context, membershipID int64, newBalance int64) error {
	query := `UPDATE customer_memberships SET pool_balance = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, newBalance, membershipID)
	if err != nil {
		return fmt.Errorf("update customer membership pool balance: %w", err)
	}
	return nil
}
