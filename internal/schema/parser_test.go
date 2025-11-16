package schema

import (
	"testing"
)

// Test_Parse_SingleFile tests parsing a single schema file (existing functionality).
func Test_Parse_SingleFile(t *testing.T) {
	parser := NewParser()

	fields, err := parser.Parse("../../testdata/simple/schema.graphql")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// simple/schema.graphql has:
	// - Query: user, users, __schema, __type (4 fields, including introspection)
	// - Mutation: createUser (1 field)
	// Total: 5 fields
	expectedCount := 5
	if len(fields) != expectedCount {
		t.Errorf("expected %d fields, got %d", expectedCount, len(fields))
	}

	// Verify field names
	fieldNames := make(map[string]bool)
	for _, f := range fields {
		fieldNames[f.FullName] = true
	}

	expected := []string{"Query.user", "Query.users", "Mutation.createUser"}
	for _, name := range expected {
		if !fieldNames[name] {
			t.Errorf("expected field %s not found", name)
		}
	}
}

// Test_Parse_Directory tests parsing multiple schema files from a directory.
func Test_Parse_Directory(t *testing.T) {
	parser := NewParser()

	fields, err := parser.Parse("../../testdata/multi-schema")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// multi-schema has:
	// - Query: user, users, post, __schema, __type (5 fields, including introspection)
	// - Mutation: createUser, updateUser (2 fields)
	// Total: 7 fields
	expectedCount := 7
	if len(fields) != expectedCount {
		t.Errorf("expected %d fields, got %d", expectedCount, len(fields))
	}

	// Verify field names
	fieldNames := make(map[string]bool)
	for _, f := range fields {
		fieldNames[f.FullName] = true
	}

	expected := []string{
		"Query.user",
		"Query.users",
		"Query.post",
		"Mutation.createUser",
		"Mutation.updateUser",
	}

	for _, name := range expected {
		if !fieldNames[name] {
			t.Errorf("expected field %s not found", name)
		}
	}
}

// Test_Parse_DirectoryMergesSchemas verifies that multiple files are merged correctly.
func Test_Parse_DirectoryMergesSchemas(t *testing.T) {
	parser := NewParser()

	fields, err := parser.Parse("../../testdata/multi-schema")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Count by type
	queryCount := 0
	mutationCount := 0

	for _, f := range fields {
		switch f.TypeName {
		case "Query":
			queryCount++
		case "Mutation":
			mutationCount++
		}
	}

	// query.graphql defines 3 Query fields + 2 introspection fields = 5 total
	if queryCount != 5 {
		t.Errorf("expected 5 Query fields (including introspection), got %d", queryCount)
	}

	// mutation.graphql defines 2 Mutation fields
	if mutationCount != 2 {
		t.Errorf("expected 2 Mutation fields, got %d", mutationCount)
	}
}

// Test_Parse_DirectoryNotFound tests error handling for non-existent directory.
func Test_Parse_DirectoryNotFound(t *testing.T) {
	parser := NewParser()

	_, err := parser.Parse("../../testdata/nonexistent")
	if err == nil {
		t.Error("expected error for non-existent directory, got nil")
	}
}

// Test_Parse_EmptyDirectory tests error handling for directory with no .graphql files.
func Test_Parse_EmptyDirectory(t *testing.T) {
	parser := NewParser()

	// testdata/simple/resolvers has .go files but no .graphql files
	_, err := parser.Parse("../../testdata/simple/resolvers")
	if err == nil {
		t.Error("expected error for directory with no .graphql files, got nil")
	}
}
