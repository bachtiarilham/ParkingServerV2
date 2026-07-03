//interface UserRepository

package user

import "context"

type Repository interface {
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id int64) (*User, error)
	UpdatePassword(ctx context.Context, userID int64, newPasswordHash string) error
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}
