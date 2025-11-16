package analyzer

import (
	"fmt"

	"github.com/nnnkkk7/graphql-unused-resolver/internal/resolver"
	"github.com/nnnkkk7/graphql-unused-resolver/internal/schema"
)

// Config contains analyzer configuration.
type Config struct {
	SchemaPath  string
	ResolverDir string
}

// Result contains analysis results.
type Result struct {
	UnusedResolvers []resolver.Method
	TotalResolvers  int
	TotalFields     int
}

// Analyzer is the main analyzer.
type Analyzer struct {
	config Config
}

// New creates a new Analyzer.
func New(config Config) *Analyzer {
	return &Analyzer{config: config}
}

// Analyze performs the complete analysis.
func (a *Analyzer) Analyze() (*Result, error) {
	// 1. Parse GraphQL schema
	schemaParser := schema.NewParser()
	fields, err := schemaParser.Parse(a.config.SchemaPath)
	if err != nil {
		return nil, fmt.Errorf("schema parse error: %w", err)
	}

	// 2. Analyze resolver code
	resolverAnalyzer := resolver.NewAnalyzer()
	resolvers, err := resolverAnalyzer.Analyze(a.config.ResolverDir)
	if err != nil {
		return nil, fmt.Errorf("resolver analysis error: %w", err)
	}

	// 3. Detect unused resolvers
	unused := a.detectUnused(fields, resolvers)

	// Build result
	result := &Result{
		UnusedResolvers: unused,
		TotalResolvers:  len(resolvers),
		TotalFields:     len(fields),
	}

	return result, nil
}

// detectUnused finds resolvers that are not defined in the schema.
func (a *Analyzer) detectUnused(fields []schema.Field, resolvers []resolver.Method) []resolver.Method {
	// Create a map of schema fields for fast lookup
	schemaMap := make(map[string]bool)
	for _, f := range fields {
		schemaMap[f.FullName] = true
	}

	// Find resolvers not in schema
	var unused []resolver.Method
	for _, r := range resolvers {
		if !schemaMap[r.GraphQLName] {
			unused = append(unused, r)
		}
	}

	return unused
}
