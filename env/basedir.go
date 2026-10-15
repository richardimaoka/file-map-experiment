package env

import "os"

func AppBaseDir() string {
	appBaseDir := os.Getenv("APP_BASE_DIR")

	if appBaseDir == "" {
		panic("APP_BASE_DIR is not set")
	}

	return appBaseDir
}
