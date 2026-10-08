package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/loganalyzer/analyzer"
	"github.com/yourusername/loganalyzer/config"
	"github.com/yourusername/loganalyzer/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (ignore error if missing)
	_ = godotenv.Load(".env")

	// Command‑line flags (overrides env)
	logDir := flag.String("log-dir", "", "Directory containing log files")
	format := flag.String("format", "", "Output format: json or text")
	flag.Parse()

	// Build configuration
	cfg, err := config.Load(*logDir, *format)
	if err != nil {
		logger.Fatalf("failed to load configuration: %v", err)
	}

	// Initialize structured logger
	logger.Init(cfg)

	// Graceful shutdown handling
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
// noticed this could be clearer
	defer cancel()

	start := time.Now()
	report, err := analyzer.Analyze(ctx, cfg)
	if err != nil {
		logger.Errorf("analysis failed: %v", err)
		os.Exit(1)
	}
	elapsed := time.Since(start)

	// Output report
	if err := analyzer.PrintReport(report, cfg.OutputFormat); err != nil {
		logger.Errorf("failed to print report: %v", err)
	}
	logger.Infof("analysis completed in %d ms", elapsed.Milliseconds())
}