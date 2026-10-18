package excel

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

// TestMain loads APP_BASE_DIR.
//
// Local: load from .env.
// If .env is not present, test fails.
// APP_BASE_DIR is usually expected in .env only, not explicitly set in terminal.
//
// CI, cloud: expect APP_BASE_DIR to be set.
// Cloud environments usually do not have .env.
func TestMain(m *testing.M) {
	if err := godotenv.Load(filepath.Join("..", ".env")); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "test setup: %v\n", err)
		os.Exit(1)
	}
	if os.Getenv("APP_BASE_DIR") == "" {
		fmt.Fprintf(os.Stderr, "test setup: APP_BASE_DIR is not set (copy .env.example to .env at repo root)\n")
		os.Exit(1)
	}
	os.Exit(m.Run())
}
