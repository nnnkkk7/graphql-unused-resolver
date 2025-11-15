package resolver

import "go/ast"

// Method represents a resolver method found in the code
type Method struct {
	// ReceiverType is the resolver type (e.g., "*queryResolver")
	ReceiverType string

	// MethodName is the Go method name (e.g., "User")
	MethodName string

	// GraphQLName is the inferred GraphQL field name (e.g., "Query.user")
	GraphQLName string

	// FilePath is the source file path
	FilePath string

	// Line is the line number in the source file
	Line int

	// ASTNode is the AST function declaration
	ASTNode *ast.FuncDecl
}
