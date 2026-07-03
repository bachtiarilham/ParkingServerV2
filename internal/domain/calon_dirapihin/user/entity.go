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
	Role             int64
	Username         string
	EmploymentStatus string
	IsVerified       bool
	RegisteredAt     time.Time
	PasswordHash     string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Record struct {
	ID           int64
	FullName     string
	PhoneNumber  string
	Email        string
	Username     string
	RoleCode     string
	PasswordHash string
	IsVerified   bool
}
