//struct User, value object

package user

import "time"

type User struct {
	ID               int64
	FullName         string
	RoleID           int64
	Nik              string
	Email            string
	Phone            string
	Role             string
	Username         string
	EmploymentStatus string
	IsVerified       bool
	RegisteredAt     time.Time
	PasswordHash     string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
