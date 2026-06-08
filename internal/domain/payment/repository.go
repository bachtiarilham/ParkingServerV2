package payment

import "context"

type Repository interface {
	Create(
		ctx context.Context,
		payment *Payment,
	) error
}
