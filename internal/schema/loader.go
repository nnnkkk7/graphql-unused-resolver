package schema

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

// loadSources loads GraphQL schema sources from a file or directory.
// If path is a file, it returns a single source.
// If path is a directory, it loads all .graphql files in the directory.
func loadSources(path string) ([]*ast.Source, error) {
	// Check if path exists
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to access path: %w", err)
	}

	// If it's a file, load it directly
	if !info.IsDir() {
		source, err := loadSingleFile(path)
		if err != nil {
			return nil, err
		}
		return []*ast.Source{source}, nil
	}

	// If it's a directory, load all .graphql files
	return loadDirectory(path)
}

// loadSingleFile loads a single GraphQL schema file.
func loadSingleFile(path string) (*ast.Source, error) {
	content, err := os.ReadFile(path) // #nosec G304 -- path is validated by caller
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return &ast.Source{
		Name:  path,
		Input: string(content),
	}, nil
}

// loadDirectory loads all .graphql files from a directory.
func loadDirectory(dirPath string) ([]*ast.Source, error) {
	var sources []*ast.Source

	// Read directory entries
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	// Filter and load .graphql files
	for _, entry := range entries {
		// Skip directories
		if entry.IsDir() {
			continue
		}

		// Only process .graphql files
		if !strings.HasSuffix(entry.Name(), ".graphql") {
			continue
		}

		// Load the file
		filePath := filepath.Join(dirPath, entry.Name())
		source, err := loadSingleFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to load %s: %w", entry.Name(), err)
		}

		sources = append(sources, source)
	}

	// Check if we found any .graphql files
	if len(sources) == 0 {
		return nil, fmt.Errorf("no .graphql files found in directory: %s", dirPath)
	}

	return sources, nil
}
