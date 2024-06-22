package domain

type Role string

const (
	UNDEFINED Role = "undefined"
	USER      Role = "user"
	ADMIN     Role = "admin"
)

type Plan string

const (
	FREE    Plan = "free"
	PREMIUM Role = "premium"
)

// User represents a user of the system.
type User struct {
	Username       string
	Password       string
	Role           Role
	Plan           Plan
	LinksRemaining int64
}

// NewUser creates a new user with the given username, password, role, plan, and links remaining.
func NewUser(username, password string, role Role, plan Plan, linksRemaining int64) *User {
	return &User{
		Username:       username,
		Password:       password,
		Role:           role,
		Plan:           plan,
		LinksRemaining: linksRemaining,
	}
}
