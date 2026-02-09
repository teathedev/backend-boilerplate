package schema

import (
	"context"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// SoftDeleteMixin implements the soft-delete pattern using 'deleted_at'.
type SoftDeleteMixin struct {
	mixin.Schema
}

// Fields of the SoftDeleteMixin.
func (SoftDeleteMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("deleted_at").
			Optional().
			Nillable(),
	}
}

type softDeleteKey struct{}

// SkipSoftDelete returns a new context that skips the soft-delete interceptor.
func SkipSoftDelete(ctx context.Context) context.Context {
	return context.WithValue(ctx, softDeleteKey{}, true)
}

// Interceptors of the SoftDeleteMixin.
func (d SoftDeleteMixin) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		ent.TraverseFunc(func(ctx context.Context, q ent.Query) error {
			// Check if we should skip the filter
			if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
				return nil
			}

			type query interface {
				WhereP(...func(*sql.Selector))
			}
			if w, ok := q.(query); ok {
				d.P(w)
			}

			return nil
		}),
	}
}

// Hooks of the SoftDeleteMixin.
func (d SoftDeleteMixin) Hooks() []ent.Hook {
	return []ent.Hook{
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				// Only intercept Delete operations
				if !m.Op().Is(ent.OpDelete | ent.OpDeleteOne) {
					return next.Mutate(ctx, m)
				}

				// We define an interface for what we need:
				// 1. Ability to set the deleted_at field
				// 2. Ability to change the Operation type
				type SoftDeleter interface {
					SetDeletedAt(time.Time)
					SetOp(ent.Op)
				}

				if sd, ok := m.(SoftDeleter); ok {
					sd.SetDeletedAt(time.Now())
					sd.SetOp(ent.OpUpdateOne)
					return next.Mutate(ctx, m)
				}

				// If the mutation doesn't support soft delete, proceed with normal delete
				return next.Mutate(ctx, m)
			})
		},
	}
}

// P adds the storage-specific "IS NULL" predicate.
func (d SoftDeleteMixin) P(w interface{ WhereP(...func(*sql.Selector)) }) {
	w.WhereP(sql.FieldIsNull("deleted_at"))
}
