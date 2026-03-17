package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/teathedev/backend-boilerplate/types"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.String("phone_number").
			MaxLen(14).
			Unique().
			NotEmpty(),
		field.String("email").
			MaxLen(255).
			Unique().
			NotEmpty(),
		field.String("username").
			MaxLen(255).
			Unique().
			NotEmpty(),
		field.Int16("role").
			GoType(types.UserRoles(0)),
		field.Int16("state").
			GoType(types.UserStates(0)),
		field.String("first_name").
			MaxLen(255).
			NotEmpty(),
		field.String("last_name").
			MaxLen(255).
			NotEmpty(),
		field.String("password_salt").
			MaxLen(36).
			Sensitive().
			NotEmpty(),
		field.String("password_hash").
			MaxLen(255).
			Sensitive().
			NotEmpty(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("refresh_tokens", RefreshToken.Type),
	}
}

// Mixins of the User.
func (User) Mixins() []ent.Mixin {
	return []ent.Mixin{
		SoftDeleteMixin{},
	}
}
