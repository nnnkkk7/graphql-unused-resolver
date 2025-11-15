package schema

import (
	"fmt"
	"os"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

// Parser parses GraphQL schema files
type Parser struct {
	schema *ast.Schema
}

// NewParser creates a new schema parser
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses a GraphQL schema file and returns fields
func (p *Parser) Parse(schemaPath string) ([]Field, error) {
	// Read schema file
	content, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema: %w", err)
	}

	// Parse with gqlparser
	source := &ast.Source{
		Name:  schemaPath,
		Input: string(content),
	}

	schema, gqlErr := gqlparser.LoadSchema(source)
	if gqlErr != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", gqlErr)
	}

	p.schema = schema

	// Extract fields from Query and Mutation types only (simplified)
	return p.extractFields()
}

// extractFields extracts all fields from Query and Mutation types
func (p *Parser) extractFields() ([]Field, error) {
	var fields []Field

	// Extract Query fields
	if queryType := p.schema.Query; queryType != nil {
		fields = append(fields, p.extractTypeFields("Query", queryType)...)
	}

	// Extract Mutation fields
	if mutationType := p.schema.Mutation; mutationType != nil {
		fields = append(fields, p.extractTypeFields("Mutation", mutationType)...)
	}

	return fields, nil
}

// extractTypeFields extracts fields from a specific type
func (p *Parser) extractTypeFields(typeName string, typeDef *ast.Definition) []Field {
	var fields []Field

	for _, field := range typeDef.Fields {
		f := Field{
			TypeName:  typeName,
			FieldName: field.Name,
			FullName:  fmt.Sprintf("%s.%s", typeName, field.Name),
		}

		fields = append(fields, f)
	}

	return fields
}
