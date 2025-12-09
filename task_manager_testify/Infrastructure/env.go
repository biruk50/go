package Infrastructure

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadEnv loads .env if present. It attempts to load a .env in the
// current directory; if not found, it walks up parent directories
// looking for a .env file so tests run from sub-packages can still
// pick up a repository-level .env.
func LoadEnv() error {
	// try default load first (current dir)
	if err := godotenv.Load(); err == nil {
		return nil
	}

	// walk up from working dir to root looking for .env
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	dir := cwd
	for {
		envPath := filepath.Join(dir, ".env")
		if _, statErr := os.Stat(envPath); statErr == nil {
			return godotenv.Load(envPath)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	// no .env found; return original Load() error for transparency
	return godotenv.Load()
}

// helper gets env with default
func GetEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
