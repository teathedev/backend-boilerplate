package filters

import (
	"github.com/teathedev/backend-boilerplate/internal/ent/accesstokenkey"
	"github.com/teathedev/backend-boilerplate/internal/ent/predicate"
	"github.com/teathedev/backend-boilerplate/types"
)

func SigningTokens() predicate.AccessTokenKey {
	return accesstokenkey.State(types.AccessTokenKeyStatesActive)
}

func VerifyTokens() predicate.AccessTokenKey {
	return accesstokenkey.StateNEQ(types.AccessTokenKeyStatesRedired)
}
