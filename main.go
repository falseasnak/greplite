package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Search options
	caseInsensitive bool
	invertMatch     bool
	showLineNumbers bool
	maxCount        int

	// Format options
	formatJSON   bool
	formatLogfmt bool
	jsonField    string

	// Output options
	colorOutput string // auto, always, never
	quietMode   bool
)

var rootCmd = &cobra.Command{
	Use:   "greplite [OPTIONS] PATTERN [FILE...]",
	Short: "A fast grep-like CLI with structured log format support",
	Long: `greplite is a fast, grep-compatible search tool with built-in support
for structured log formats like JSON and logfmt.

Examples:
  greplite "error" app.log
  greplite --json --field message "timeout" app.log
  greplite --logfmt --field level=error "" app.log
  cat app.log | greplite "panic"`,
	Args:         cobra.MinimumNArgs(1),
	RunE:         runSearch,
	SilenceUsage: true,
}

func init() {
	// Search flags
	rootCmd.Flags().BoolVarP(&caseInsensitive, "ignore-case", "i", false, "Case-insensitive matching")
	rootCmd.Flags().BoolVarP(&invertMatch, "invert-match", "v", false, "Select non-matching lines")
	rootCmd.Flags().BoolVarP(&showLineNumbers, "line-number", "n", false, "Print line numbers")
	rootCmd.Flags().IntVarP(&maxCount, "max-count", "m", 0, "Stop after NUM matches (0 = unlimited)")

	// Format flags
	rootCmd.Flags().BoolVar(&formatJSON, "json", false, "Parse input as JSON log lines")
	rootCmd.Flags().BoolVar(&formatLogfmt, "logfmt", false, "Parse input as logfmt log lines")
	rootCmd.Flags().StringVar(&jsonField, "field", "", "Search within a specific field (use with --json or --logfmt)")

	// Output flags
	rootCmd.Flags().StringVar(&colorOutput, "color", "auto", "Colorize output: auto, always, never")
	rootCmd.Flags().BoolVarP(&quietMode, "quiet", "q", false, "Suppress output, exit 0 if match found")
}

func runSearch(cmd *cobra.Command, args []string) error {
	pattern := args[0]
	files := args[1:]

	// Build search configuration
	cfg := &SearchConfig{
		Pattern:         pattern,
		CaseInsensitive: caseInsensitive,
		InvertMatch:     invertMatch,
		ShowLineNumbers: showLineNumbers,
		MaxCount:        maxCount,
		FormatJSON:      formatJSON,
		FormatLogfmt:    formatLogfmt,
		Field:           jsonField,
		ColorOutput:     colorOutput,
		QuietMode:       quietMode,
	}

	searcher, err := NewSearcher(cfg)
	if err != nil {
		return fmt.Errorf("invalid pattern: %w", err)
	}

	// Search stdin if no files provided
	if len(files) == 0 {
		matched, err := searcher.SearchReader(os.Stdin, "")
		if err != nil {
			return err
		}
		if !matched {
			os.Exit(1)
		}
		return nil
	}

	// Search each file
	anyMatched := false
	showFilename := len(files) > 1

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "greplite: %s: %v\n", file, err)
			continue
		}

		filename := ""
		if showFilename {
			filename = file
		}

		matched, err := searcher.SearchReader(f, filename)
		f.Close()

		if err != nil {
			fmt.Fprintf(os.Stderr, "greplite: %s: %v\n", file, err)
			continue
		}

		if matched {
			anyMatched = true
		}
	}

	if !anyMatched {
		os.Exit(1)
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(2)
	}
}
