package resolvers

import (
	"context"
)

// Resolver is the root resolver
type Resolver struct{}

// queryResolver implements the Query type
type queryResolver struct{ *Resolver }

// mutationResolver implements the Mutation type
type mutationResolver struct{ *Resolver }

// Query returns the query resolver
func (r *Resolver) Query() *queryResolver {
	return &queryResolver{r}
}

// Mutation returns the mutation resolver
func (r *Resolver) Mutation() *mutationResolver {
	return &mutationResolver{r}
}

// User resolver (defined in schema)
func (r *queryResolver) User(ctx context.Context, id string) (*User, error) {
	return &User{ID: id, Name: "Test User"}, nil
}

// Users resolver (defined in schema)
func (r *queryResolver) Users(ctx context.Context) ([]*User, error) {
	return []*User{{ID: "1", Name: "User 1"}}, nil
}

// CreateUser resolver (defined in schema)
func (r *mutationResolver) CreateUser(ctx context.Context, name string) (*User, error) {
	return &User{ID: "new", Name: name}, nil
}

// ============================================
// UNUSED RESOLVERS (not in schema)
// ============================================

// Orders resolver (NOT in schema - should be detected as unused)
func (r *queryResolver) Orders(ctx context.Context) ([]*Order, error) {
	return []*Order{{ID: "1"}}, nil
}

// LegacyField resolver (NOT in schema - should be detected as unused)
func (r *queryResolver) LegacyField(ctx context.Context) (string, error) {
	return "legacy", nil
}

// DeleteUser resolver (NOT in schema - should be detected as unused)
func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	return true, nil
}

// User model
type User struct {
	ID    string
	Name  string
	Email string
}

// Order model (not in schema)
type Order struct {
	ID string
}
