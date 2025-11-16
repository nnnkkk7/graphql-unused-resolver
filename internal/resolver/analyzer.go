package resolver

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// Analyzer analyzes resolver code.
type Analyzer struct {
	fset *token.FileSet
}

// NewAnalyzer creates a new resolver analyzer.
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		fset: token.NewFileSet(),
	}
}

// Analyze analyzes a resolver directory and returns all resolver methods.
func (a *Analyzer) Analyze(dir string) ([]Method, error) {
	var methods []Method
	var errors []error

	// Walk through .go files in the directory
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-Go files
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Skip test files
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Parse the file
		fileMethods, fileErrs := a.analyzeFile(path)
		methods = append(methods, fileMethods...)
		if len(fileErrs) > 0 {
			errors = append(errors, fileErrs...)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	if len(errors) > 0 {
		// Return first error for simplicity
		return methods, errors[0]
	}

	return methods, nil
}

// analyzeFile analyzes a single Go file.
func (a *Analyzer) analyzeFile(path string) ([]Method, []error) {
	var methods []Method
	var errors []error

	// Parse the file
	file, err := parser.ParseFile(a.fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to parse %s: %w", path, err)}
	}

	// Walk through function declarations
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// Skip functions without receivers (not methods)
		if funcDecl.Recv == nil {
			continue
		}

		// Try to extract resolver method
		method, err := a.extractResolverMethod(funcDecl, path)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		if method != nil {
			methods = append(methods, *method)
		}
	}

	return methods, errors
}

// extractResolverMethod extracts resolver method information from a function declaration.
func (a *Analyzer) extractResolverMethod(funcDecl *ast.FuncDecl, filePath string) (*Method, error) {
	// Get receiver type
	receiverType := a.getReceiverType(funcDecl.Recv)
	if receiverType == "" {
		return nil, nil
	}

	// Check if it's a resolver (gqlgen pattern: *queryResolver, *mutationResolver, etc.)
	if !strings.HasSuffix(strings.ToLower(receiverType), "resolver") {
		return nil, nil
	}

	// Infer GraphQL type from receiver type
	// *queryResolver -> Query
	// *mutationResolver -> Mutation
	graphqlType := a.inferGraphQLType(receiverType)
	methodName := funcDecl.Name.Name

	// Infer GraphQL field name (lowercase first letter)
	fieldName := lowerFirst(methodName)
	fullName := fmt.Sprintf("%s.%s", graphqlType, fieldName)

	// Get position
	pos := a.fset.Position(funcDecl.Pos())

	return &Method{
		ReceiverType: receiverType,
		MethodName:   methodName,
		GraphQLName:  fullName,
		FilePath:     filePath,
		Line:         pos.Line,
		ASTNode:      funcDecl,
	}, nil
}

// getReceiverType extracts the receiver type from a method.
func (a *Analyzer) getReceiverType(recv *ast.FieldList) string {
	if recv == nil || len(recv.List) == 0 {
		return ""
	}

	field := recv.List[0]

	switch typ := field.Type.(type) {
	case *ast.StarExpr:
		// *queryResolver
		if ident, ok := typ.X.(*ast.Ident); ok {
			return "*" + ident.Name
		}
	case *ast.Ident:
		// queryResolver
		return typ.Name
	}

	return ""
}

// inferGraphQLType infers GraphQL type from receiver type.
// Examples:
// *queryResolver -> Query
// *mutationResolver -> Mutation
// *userResolver -> User.
func (a *Analyzer) inferGraphQLType(receiverType string) string {
	// Remove pointer prefix
	typeName := strings.TrimPrefix(receiverType, "*")

	// Remove "Resolver" suffix
	typeName = strings.TrimSuffix(typeName, "Resolver")

	// Uppercase first letter
	return upperFirst(typeName)
}

// Helper functions

// lowerFirst converts the first character to lowercase.
func lowerFirst(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// upperFirst converts the first character to uppercase.
func upperFirst(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
