package schema

// Field represents a GraphQL field defined in the schema
type Field struct {
	// TypeName is the parent type (e.g., "Query", "Mutation")
	TypeName string

	// FieldName is the field name (e.g., "user", "createOrder")
	FieldName string

	// FullName is the fully qualified name (e.g., "Query.user")
	FullName string
}
