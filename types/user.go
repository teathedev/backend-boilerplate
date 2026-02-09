package types

import (
	"time"

	"github.com/google/uuid"
)

// UserStates: Begin

type UserStates int16

const (
	UserStatesCreated UserStates = 0
	UserStatesInvited UserStates = 1
	UserStatesActive  UserStates = 2
	UserStatesPassive UserStates = 3
	UserStatesBanned  UserStates = 4
)

func (UserStates) Values() []int16 {
	return []int16{
		int16(UserStatesCreated),
		int16(UserStatesInvited),
		int16(UserStatesActive),
		int16(UserStatesPassive),
		int16(UserStatesBanned),
	}
}

// UserStates: End

// UserRoles: Begin

type UserRoles int16

const (
	UserRolesSuperUser  UserRoles = 0
	UserRolesAgent      UserRoles = 1
	UserRolesClient     UserRoles = 2
	UserRolesContractor UserRoles = 3
)

func (UserRoles) Values() []int16 {
	return []int16{
		int16(UserRolesSuperUser),
		int16(UserRolesAgent),
		int16(UserRolesClient),
		int16(UserRolesContractor),
	}
}

// UserRoles: End

type User struct {
	ID           uuid.UUID  `json:"id"`
	PhoneNumber  string     `json:"phoneNumber"`
	Email        string     `json:"email"`
	Username     string     `json:"username"`
	Role         UserRoles  `json:"role"`
	State        UserStates `json:"state"`
	FirstName    string     `json:"firstName"`
	LastName     string     `json:"lastName"`
	PasswordSalt string     `json:"-"`
	PasswordHash string     `json:"-"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `json:"-"`
}
