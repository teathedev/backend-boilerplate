package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/teathedev/backend-boilerplate/types"
)

// AccessTokenKey holds the schema definition for the AccessTokenKey entity.
type AccessTokenKey struct {
	ent.Schema
}

// Fields of the AccessTokenKey.
func (AccessTokenKey) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.Bytes("private_key_encrypted").
			NotEmpty(),
		field.String("public_pem").
			NotEmpty(),
		field.Int8("state").
			GoType(types.AccessTokenKeyStates(0)),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the AccessTokenKey.
func (AccessTokenKey) Edges() []ent.Edge {
	return nil
}

// Mixins of the AccessTokenKey.
func (AccessTokenKey) Mixins() []ent.Mixin {
	return []ent.Mixin{
		SoftDeleteMixin{},
	}
}
