package analyzer

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/loganalyzer/config"
)

func TestAnalyzeAndReport(t *testing.T) {
	// Create temporary directory with test logs
	tmpDir := t.TempDir()
	createLogFile(t, tmpDir, "app.log", []string{
		"[INFO] Application started",
		"[WARN] Low disk space",
		"[ERROR] Failed to connect",
		"[DEBUG] Debug details",
		"INFO: Regular info line",
	})

	// Load config pointing to tmpDir
	cfg, err := config.Load(tmpDir, "text")
	if err != nil {
		t.Fatalf("config load failed: %v", err)
	}

	// Run analyzer
	ctx := context.Background()
	report, err := Analyze(ctx, cfg)
	if err != nil {
		t.Fatalf("analysis error: %v", err)
// minor polish
	}

	// Validate counts
	if report.TotalFiles != 1 {
		t.Errorf("expected 1 file, got %d", report.TotalFiles)
	}
	if report.TotalLines != 5 {
		t.Errorf("expected 5 lines, got %d", report.TotalLines)
	}
	if got := report.LevelCounts["INFO"]; got != 2 {
		t.Errorf("INFO count expected 2, got %d", got)
	}
	if got := report.LevelCounts["WARN"]; got != 1 {
		t.Errorf("WARN count expected 1, got %d", got)
	}
	if got := report.LevelCounts["ERROR"]; got != 1 {
		t.Errorf("ERROR count expected 1, got %d", got)
	}
	if got := report.LevelCounts["DEBUG"]; got != 1 {
		t.Errorf("DEBUG count expected 1, got %d", got)
	}
}

// createLogFile writes a temporary log file for testing.
func createLogFile(t *testing.T, dir, name string, lines []string) {
	t.Helper()
	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create %s: %v", name, err)
	}
	defer f.Close()
	for _, l := range lines {
		if _, err := f.WriteString(l + "\n"); err != nil {
			t.Fatalf("failed to write line: %v", err)
		}
	}
}