// Package output handles formatting and displaying HTTP responses.
package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/EdgarOrtegaRamirez/httpforge/models"
	"github.com/fatih/color"
)

var (
	green  = color.New(color.FgGreen)
	red    = color.New(color.FgRed)
	yellow = color.New(color.FgYellow)
	cyan   = color.New(color.FgCyan)
	bold   = color.New(color.Bold)
	dim    = color.New(color.Faint)
)

// PrintResponse displays an HTTP response in the specified format.
func PrintResponse(resp *models.Response, format models.OutputFormat, verbose bool) {
	switch format {
	case models.OutputJSON:
		printJSON(resp)
	case models.OutputCSV:
		printCSV(resp)
	default:
		printText(resp, verbose)
	}
}

// printText displays the response in a human-readable text format.
func printText(resp *models.Response, verbose bool) {
	// Status line
	statusColor := green
	if resp.StatusCode >= 400 {
		statusColor = red
	} else if resp.StatusCode >= 300 {
		statusColor = yellow
	}

	fmt.Print("\n")
	bold.Print("Response\n")
	fmt.Print(strings.Repeat("─", 50) + "\n")

	// Status
	fmt.Print("Status:    ")
	statusColor.Printf("%d", resp.StatusCode)
	fmt.Println()

	// Timing
	fmt.Printf("Time:      %s\n", resp.Timing.Round(0))
	fmt.Printf("Size:      %s\n", formatSize(resp.Size))

	if verbose {
		// Headers
		fmt.Print("\n")
		bold.Println("Headers:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		for key, value := range resp.Headers {
			fmt.Fprintf(w, "  %s:\t%s\n", dim.Sprint(key), value)
		}
		w.Flush()
	}

	// Body
	if resp.Body != "" {
		fmt.Print("\n")
		bold.Println("Body:")
		if isJSON(resp.Body) {
			prettyJSON := prettyPrintJSON(resp.Body)
			fmt.Println(prettyJSON)
		} else {
			fmt.Println(resp.Body)
		}
	}
}

// printJSON displays the response in JSON format.
func printJSON(resp *models.Response) {
	data := map[string]interface{}{
		"status_code": resp.StatusCode,
		"headers":     resp.Headers,
		"body":        resp.Body,
		"timing":      resp.Timing.String(),
		"size":        resp.Size,
	}

	if resp.ExtractedVars != nil {
		data["extracted_vars"] = resp.ExtractedVars
	}

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		return
	}
	fmt.Println(string(jsonBytes))
}

// printCSV displays the response summary in CSV format.
func printCSV(resp *models.Response) {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	// Header
	writer.Write([]string{"field", "value"})

	// Status
	writer.Write([]string{"status_code", strconv.Itoa(resp.StatusCode)})

	// Timing
	writer.Write([]string{"timing", resp.Timing.String()})

	// Size
	writer.Write([]string{"size", strconv.Itoa(resp.Size)})

	// Body length
	writer.Write([]string{"body_length", strconv.Itoa(len(resp.Body))})

	// Key headers
	for key, value := range resp.Headers {
		writer.Write([]string{"header." + key, value})
	}
}

// PrintAssertions displays assertion results.
func PrintAssertions(results []models.AssertionResult) {
	fmt.Print("\n")
	bold.Println("Assertions:")
	fmt.Print(strings.Repeat("─", 50) + "\n")

	for _, r := range results {
		if r.Passed {
			green.Printf("  ✓ %s\n", r.Message)
		} else {
			red.Printf("  ✗ %s\n", r.Message)
		}
	}

	// Summary
	total := len(results)
	passed := 0
	for _, r := range results {
		if r.Passed {
			passed++
		}
	}

	fmt.Print("\n")
	if passed == total {
		green.Printf("All %d assertions passed\n", total)
	} else {
		red.Printf("%d/%d assertions failed\n", total-passed, total)
	}
}

// PrintError displays an error message.
func PrintError(err error) {
	red.Printf("Error: %v\n", err)
}

// PrintInfo displays an informational message.
func PrintInfo(msg string) {
	cyan.Printf("ℹ %s\n", msg)
}

// PrintSuccess displays a success message.
func PrintSuccess(msg string) {
	green.Printf("✓ %s\n", msg)
}

// PrintWarning displays a warning message.
func PrintWarning(msg string) {
	yellow.Printf("⚠ %s\n", msg)
}

// isJSON checks if a string is valid JSON.
func isJSON(s string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(s), &js) == nil
}

// prettyPrintJSON pretty-prints a JSON string.
func prettyPrintJSON(s string) string {
	var js interface{}
	if err := json.Unmarshal([]byte(s), &js); err != nil {
		return s
	}
	pretty, err := json.MarshalIndent(js, "", "  ")
	if err != nil {
		return s
	}
	return string(pretty)
}

// formatSize formats a byte size in human-readable form.
func formatSize(bytes int) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
