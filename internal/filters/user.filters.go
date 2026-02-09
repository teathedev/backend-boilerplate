package filters

import (
	"github.com/teathedev/backend-boilerplate/internal/ent/predicate"
	"github.com/teathedev/backend-boilerplate/internal/ent/user"
	"github.com/teathedev/backend-boilerplate/types"
)

func ActiveUsers() predicate.User {
	return user.StateIn(types.UserStatesActive, types.UserStatesInvited)
}
