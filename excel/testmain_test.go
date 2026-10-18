package excel

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

// TestMain loads repo-root ../.env when present, then exits if APP_BASE_DIR is unset.
// Same logic for local, CI, and other hosts; only the usual source of APP_BASE_DIR differs.
//
// Local: copy .env.example to repo-root .env (absolute module root). No .env and no export → exit here.
// CI / cloud: set APP_BASE_DIR in the environment; a .env file is optional. godotenv ignores a missing file.
// godotenv never overrides APP_BASE_DIR already set in the environment.
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
