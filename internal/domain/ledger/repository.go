package ledger

import "context"

type Repository interface {
	CreateEntries(
		ctx context.Context,
		entries []Entry,
	) error
}
