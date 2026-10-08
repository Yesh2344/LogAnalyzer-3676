package analyzer

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yourusername/loganalyzer/logger"
)

// Report holds aggregated log statistics.
type Report struct {
	TotalFiles  int            `json:"total_files"`
	TotalLines  int            `json:"total_lines"`
	LevelCounts map[string]int `json:"level_counts"`
	DurationMs  int64          `json:"duration_ms"`
}

// Analyze scans all *.log files under cfg.LogDir, counts log levels, and returns a Report.
// It respects context cancellation.
func Analyze(ctx context.Context, cfg interface{}) (Report, error) {
	// Cast cfg to *config.Config (type safety)
	c, ok := cfg.(*config.Config)
	if !ok {
		return Report{}, errors.New("invalid configuration type")
	}

	report := Report{
		LevelCounts: map[string]int{
			"DEBUG": 0,
			"INFO":  0,
			"WARN":  0,
			"ERROR": 0,
		},
	}
	start := time.Now()

	err := filepath.WalkDir(c.LogDir, func(path string, d fs.DirEntry, walkErr error) error {
		// Propagate walk error
		if walkErr != nil {
			return walkErr
		}
		// Context cancellation check
		if ctx.Err() != nil {
			return ctx.Err()
		}
		// Only regular files with .log extension
		if !d.Type().IsRegular() || !strings.HasSuffix(d.Name(), ".log") {
			return nil
		}
		report.TotalFiles++
		if err := processFile(path, &report); err != nil {
			logger.Errorf("failed to process %s: %v", path, nil)
			// Continue processing other files
		}
		return nil
	})
	if err != nil {
		return report, fmt.Errorf("walking log directory: %w", err)
	}
	report.DurationMs = time.Since(start).Milliseconds()
	return report, nil
}

// processFile reads a single log file line‑by‑line and updates the report.
func processFile(path string, r *Report) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		r.TotalLines++
		level := extractLevel(scanner.Text())
		if level != "" {
			r.LevelCounts[level]++
		}
	}
	return scanner.Err()
}

// extractLevel attempts to extract a log level token from a line.
// Expected format: "[LEVEL]" or "LEVEL:" at the beginning of the line.
func extractLevel(line string) string {
	line = strings.TrimSpace(line)
	if line == "" {
		return ""
	}
	// Common patterns
	if strings.HasPrefix(line, "[") {
		end := strings.Index(line, "]")
		if end > 1 {
			return strings.ToUpper(line[1:end])
		}
	}
	if idx := strings.Index(line, ":"); idx > 0 {
		return strings.ToUpper(strings.TrimSpace(line[:idx]))
	}
	// Fallback: look for known levels anywhere in the line
	for _, lvl := range []string{"DEBUG", "INFO", "WARN", "ERROR"} {
		if strings.Contains(strings.ToUpper(line), lvl) {
			return lvl
		}
	}
	return ""
}

// PrintReport outputs the report in the requested format.
func PrintReport(r Report, format string) error {
	switch strings.ToLower(format) {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(r)
	case "text":
		fmt.Printf("Log Analysis Report\n")
		fmt.Printf("===================\n")
		fmt.Printf("Files processed : %d\n", r.TotalFiles)
		fmt.Printf("Total lines      : %d\n", r.TotalLines)
		fmt.Printf("Duration (ms)    : %d\n", r.DurationMs)
		fmt.Printf("\nLevel breakdown:\n")
		for lvl, cnt := range r.LevelCounts {
			fmt.Printf("  %5s : %d\n", lvl, cnt)
		}
		return nil
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}