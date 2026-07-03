package topup

import "context"

type Repository interface {
	Create(
		ctx context.Context,
		req *Request,
	) error

	UpdateStatus(
		ctx context.Context,
		id int64,
		status string,
	) error
}
