# graphql-unused-resolver

Detect unused GraphQL resolvers and their dependencies (services, repositories, etc.) by analyzing your schema and Go backend code.

## Features

- ✅ **GraphQL Schema Analysis**: Parse your GraphQL schema to extract all defined fields
- ✅ **Resolver Detection**: Analyze Go code to find all gqlgen-style resolver implementations
- ✅ **Unused Resolver Detection**: Identify resolvers that exist in code but not in the schema
- ✅ **Simple CLI**: Easy-to-use command-line interface
- 📊 **Clear Reports**: Human-readable output with file locations

## Installation

### From Source

```bash
git clone https://github.com/s20590/graphql-unused-resolver.git
cd graphql-unused-resolver
go build -o bin/graphql-unused-resolver ./cmd/graphql-unused-resolver
```

### Using Go Install

```bash
go install github.com/s20590/graphql-unused-resolver/cmd/graphql-unused-resolver@latest
```

## Usage

### Basic Usage

```bash
graphql-unused-resolver --schema schema.graphql --resolvers ./graph/resolvers
```

### Example Output

```
==========================================
GraphQL Unused Resolver Analysis Report
==========================================

Total Schema Fields: 5
Total Resolvers:     8
Unused Resolvers:    3

❌ Unused Resolvers (3):

1. Query.orders
   Receiver: *queryResolver
   Method:   Orders
   Location: resolvers/query.go:46

2. Query.legacyField
   Receiver: *queryResolver
   Method:   LegacyField
   Location: resolvers/query.go:51

3. Mutation.deleteUser
   Receiver: *mutationResolver
   Method:   DeleteUser
   Location: resolvers/mutation.go:56

==========================================
💡 Recommendation: Review and remove 3 unused resolver(s)
==========================================
```

## How It Works

1. **Schema Parsing**: Reads your GraphQL schema and extracts all Query and Mutation fields
2. **Resolver Analysis**: Scans your Go code for gqlgen-style resolver methods (`*queryResolver`, `*mutationResolver`)
3. **Comparison**: Compares schema fields with implemented resolvers
4. **Detection**: Identifies resolvers that are implemented but not defined in the schema

## Supported Patterns

Currently supports **gqlgen** resolver patterns:

```go
type queryResolver struct{ *Resolver }

// This resolver will be detected as "Query.user"
func (r *queryResolver) User(ctx context.Context, id string) (*User, error) {
    // ...
}
```

## Project Structure

```
.
├── cmd/
│   └── graphql-unused-resolver/
│       └── main.go              # CLI entry point
├── internal/
│   ├── schema/
│   │   ├── parser.go            # GraphQL schema parser
│   │   └── types.go             # Schema types
│   └── resolver/
│       ├── analyzer.go          # Resolver code analyzer
│       └── types.go             # Resolver types
├── pkg/
│   └── analyzer/
│       └── analyzer.go          # Main analysis logic
├── testdata/
│   └── simple/                  # Test data
└── README.md
```

## Minimum Viable Product (MVP)

This is the **minimal MVP** version with core functionality:

### ✅ Implemented
- GraphQL schema parsing (Query/Mutation types only)
- gqlgen resolver detection
- Unused resolver detection
- Simple text output

### 🔜 Future Features (Phase 2+)
- Dependency tracking (services, repositories)
- Confidence scoring
- Multiple output formats (JSON, Markdown)
- Whitelist/ignore patterns
- Configuration file support
- Support for custom GraphQL types
- Subscription type support

## Development

### Build

```bash
make build
# or
go build -o bin/graphql-unused-resolver ./cmd/graphql-unused-resolver
```

### Test

```bash
# Run with test data
./bin/graphql-unused-resolver \
  --schema testdata/simple/schema.graphql \
  --resolvers testdata/simple/resolvers
```

### Requirements

- Go 1.21 or higher
- GraphQL schema file
- Go resolver code following gqlgen patterns

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) file for details

## Acknowledgments

- Built with [gqlparser](https://github.com/vektah/gqlparser) for GraphQL schema parsing
- Uses Go's standard `go/ast` for code analysis
- CLI powered by [Cobra](https://github.com/spf13/cobra)
