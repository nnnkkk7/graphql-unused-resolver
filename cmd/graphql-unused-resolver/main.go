package main

import (
	"fmt"
	"os"

	"github.com/nnnkkk7/graphql-unused-resolver/pkg/analyzer"

	"github.com/spf13/cobra"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	var config analyzer.Config

	cmd := &cobra.Command{
		Use:   "graphql-unused-resolver",
		Short: "Detect unused GraphQL resolvers and their dependencies",
		Long: `Detect unused GraphQL resolvers and their dependencies by analyzing
your schema and Go backend code.

This tool helps you identify resolver code that is no longer needed
because the corresponding fields have been removed from your GraphQL schema.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAnalysis(config)
		},
	}

	// Define flags
	cmd.Flags().StringVar(&config.SchemaPath, "schema", "", "Path to GraphQL schema file (required)")
	cmd.Flags().StringVar(&config.ResolverDir, "resolvers", "", "Path to resolver directory (required)")

	_ = cmd.MarkFlagRequired("schema")
	_ = cmd.MarkFlagRequired("resolvers")

	return cmd
}

func runAnalysis(config analyzer.Config) error {
	// Validate configuration
	if err := validateConfig(config); err != nil {
		return err
	}

	// Create analyzer
	a := analyzer.New(config)

	// Run analysis
	result, err := a.Analyze()
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	// Print report
	printReport(result)

	return nil
}

func validateConfig(config analyzer.Config) error {
	// Check schema file exists
	if _, err := os.Stat(config.SchemaPath); os.IsNotExist(err) {
		return fmt.Errorf("schema file does not exist: %s", config.SchemaPath)
	}

	// Check resolver directory exists
	if info, err := os.Stat(config.ResolverDir); os.IsNotExist(err) {
		return fmt.Errorf("resolver directory does not exist: %s", config.ResolverDir)
	} else if !info.IsDir() {
		return fmt.Errorf("resolver path is not a directory: %s", config.ResolverDir)
	}

	return nil
}

func printReport(result *analyzer.Result) {
	fmt.Println("==========================================")
	fmt.Println("GraphQL Unused Resolver Analysis Report")
	fmt.Println("==========================================")
	fmt.Println()

	// Summary
	fmt.Printf("Total Schema Fields: %d\n", result.TotalFields)
	fmt.Printf("Total Resolvers:     %d\n", result.TotalResolvers)
	fmt.Printf("Unused Resolvers:    %d\n", len(result.UnusedResolvers))
	fmt.Println()

	// Unused resolvers
	if len(result.UnusedResolvers) == 0 {
		fmt.Println("✅ No unused resolvers found!")
		return
	}

	fmt.Printf("❌ Unused Resolvers (%d):\n", len(result.UnusedResolvers))
	fmt.Println()

	for i, r := range result.UnusedResolvers {
		fmt.Printf("%d. %s\n", i+1, r.GraphQLName)
		fmt.Printf("   Receiver: %s\n", r.ReceiverType)
		fmt.Printf("   Method:   %s\n", r.MethodName)
		fmt.Printf("   Location: %s:%d\n", r.FilePath, r.Line)
		fmt.Println()
	}

	fmt.Println("==========================================")
	fmt.Printf("💡 Recommendation: Review and remove %d unused resolver(s)\n", len(result.UnusedResolvers))
	fmt.Println("==========================================")
}
