// Package cmd implements the CLI commands for HttpForge.
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/EdgarOrtegaRamirez/httpforge/assert"
	"github.com/EdgarOrtegaRamirez/httpforge/engine"
	"github.com/EdgarOrtegaRamirez/httpforge/export"
	"github.com/EdgarOrtegaRamirez/httpforge/models"
	"github.com/EdgarOrtegaRamirez/httpforge/output"
	"github.com/EdgarOrtegaRamirez/httpforge/parser"
	"github.com/EdgarOrtegaRamirez/httpforge/template"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	outputFormat string
	verbose      bool
	envFile      string
	vars         []string
	bold         = color.New(color.Bold)
)

var rootCmd = &cobra.Command{
	Use:   "httpforge",
	Short: "HTTP API Testing Toolkit",
	Long: `HttpForge is a lightweight, config-driven HTTP API testing tool.
Build, execute, and validate HTTP requests from YAML/JSON configs.`,
}

var runCmd = &cobra.Command{
	Use:   "run <file|url>",
	Short: "Execute an HTTP request",
	Long: `Execute an HTTP request from a file or URL.
Supports YAML/JSON configs or simple URL shorthand.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRequest(args[0])
	},
}

var exportCmd = &cobra.Command{
	Use:   "export <file>",
	Short: "Export request as code",
	Long: `Export an HTTP request as curl, Python, or JavaScript code.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return exportRequest(args[0], cmd)
	},
}

var infoCmd = &cobra.Command{
	Use:   "info <file>",
	Short: "Display request configuration",
	Long: `Parse and display the request configuration without executing it.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return showInfo(args[0])
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate <file>",
	Short: "Validate request file format",
	Long: `Parse and validate a request configuration file.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return validateFile(args[0])
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("HttpForge v1.0.0")
		fmt.Println("A lightweight HTTP API testing toolkit")
	},
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text, json, csv)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().StringVarP(&envFile, "env", "e", "", "Environment file (YAML/JSON)")
	rootCmd.PersistentFlags().StringArrayVarP(&vars, "var", "V", []string{}, "Set variable (key=value)")

	// Export format flag
	exportCmd.Flags().StringP("format", "f", "curl", "Export format (curl, python, javascript)")

	// Add commands
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(versionCmd)
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// runRequest executes an HTTP request.
func runRequest(input string) error {
	// Build template engine
	eng := template.NewEngine(nil)

	// Load env file if specified
	if envFile != "" {
		env, err := parser.LoadEnvironment(envFile)
		if err != nil {
			return fmt.Errorf("loading environment: %w", err)
		}
		eng.Merge(env.Variables)
	}

	// Load OS environment variables
	eng.LoadEnv("HTTPFORGE_")

	// Parse inline variables
	inlineVars := parseVars(vars)
	eng.Merge(inlineVars)

	// Load request
	var req *models.Request
	var err error

	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		req, err = parser.ParseRequest(input)
	} else {
		config, loadErr := parser.LoadRequestFile(input)
		if loadErr != nil {
			return fmt.Errorf("loading request: %w", loadErr)
		}
		req = &config.Request

		// Merge collection variables
		if config.Environment != nil {
			eng.Merge(config.Environment.Variables)
		}
	}

	if err != nil {
		return fmt.Errorf("parsing request: %w", err)
	}

	// Validate template variables
	if missing := eng.Validate(req.URL); len(missing) > 0 {
		output.PrintWarning(fmt.Sprintf("Missing variables: %s", strings.Join(missing, ", ")))
	}

	// Render templates
	req.URL = eng.Render(req.URL)
	req.Body = eng.Render(req.Body)
	req.Headers = eng.RenderMap(req.Headers)
	req.QueryParams = eng.RenderMap(req.QueryParams)

	// Execute request
	fmt.Printf("\n%s %s\n", bold.Sprint(req.Method), req.URL)

	client := engine.NewClient()
	resp, err := client.Execute(req)
	if err != nil {
		output.PrintError(err)
		return err
	}

	// Print response
	format := models.OutputFormat(outputFormat)
	output.PrintResponse(resp, format, verbose)

	// Run assertions if any
	if len(req.Assertions) > 0 {
		results := assert.Validate(resp, req.Assertions)
		output.PrintAssertions(results)

		if !assert.AllPassed(results) {
			os.Exit(1)
		}
	}

	return nil
}

// exportRequest exports a request as code.
func exportRequest(input string, cmd *cobra.Command) error {
	config, err := parser.LoadRequestFile(input)
	if err != nil {
		return fmt.Errorf("loading request: %w", err)
	}

	formatFlag, _ := cmd.Flags().GetString("format")
	var format models.ExportFormat
	switch strings.ToLower(formatFlag) {
	case "python", "py":
		format = models.ExportPython
	case "javascript", "js":
		format = models.ExportJS
	default:
		format = models.ExportCurl
	}

	code := export.Export(&config.Request, format)
	fmt.Println(code)

	return nil
}

// showInfo displays request configuration.
func showInfo(input string) error {
	config, err := parser.LoadRequestFile(input)
	if err != nil {
		return fmt.Errorf("loading request: %w", err)
	}

	req := config.Request
	fmt.Print("\n")
	bold.Println("Request Configuration")
	fmt.Print(strings.Repeat("─", 50) + "\n")

	if req.Name != "" {
		fmt.Printf("Name:        %s\n", req.Name)
	}
	fmt.Printf("Method:      %s\n", req.Method)
	fmt.Printf("URL:         %s\n", req.URL)

	if len(req.Headers) > 0 {
		fmt.Println("\nHeaders:")
		for key, value := range req.Headers {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	if req.Body != "" {
		fmt.Printf("\nBody (%s):\n%s\n", req.BodyType, req.Body)
	}

	if len(req.QueryParams) > 0 {
		fmt.Println("\nQuery Parameters:")
		for key, value := range req.QueryParams {
			fmt.Printf("  %s=%s\n", key, value)
		}
	}

	if len(req.Assertions) > 0 {
		fmt.Println("\nAssertions:")
		for _, a := range req.Assertions {
			fmt.Printf("  %s %s %q\n", a.Field, a.Op, a.Value)
		}
	}

	// Check for template variables
	var allText strings.Builder
	allText.WriteString(req.URL)
	allText.WriteString(req.Body)
	for _, h := range req.Headers {
		allText.WriteString(h)
	}

	if template.HasVariables(allText.String()) {
		vars := template.ListVariables(allText.String())
		fmt.Printf("\nTemplate Variables: %s\n", strings.Join(vars, ", "))
	}

	return nil
}

// validateFile validates a request configuration file.
func validateFile(input string) error {
	_, err := parser.LoadRequestFile(input)
	if err != nil {
		output.PrintError(err)
		return err
	}

	output.PrintSuccess(fmt.Sprintf("Valid configuration: %s", input))
	return nil
}

// parseVars parses key=value variable strings.
func parseVars(varStrs []string) map[string]string {
	vars := make(map[string]string)
	for _, v := range varStrs {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) == 2 {
			vars[parts[0]] = parts[1]
		}
	}
	return vars
}
