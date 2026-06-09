//interface UserRepository

package user

import "context"

type Repository interface {
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}
