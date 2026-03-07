package filters

import (
	"strings"

	"github.com/teathedev/backend-boilerplate/internal/ent/predicate"
	"github.com/teathedev/backend-boilerplate/internal/ent/user"
	"github.com/teathedev/backend-boilerplate/types"
)

func ActiveUsers() predicate.User {
	return user.StateIn(types.UserStatesActive, types.UserStatesInvited)
}

func UserByIdentifier(identifier string) predicate.User {
	return user.
		Or(
			user.PhoneNumberEQ(identifier),
			user.EmailEQ(strings.ToLower(identifier)),
			user.UsernameEQ(strings.ToLower(identifier)),
		)
}
