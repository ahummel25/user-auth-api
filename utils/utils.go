package utils

import "os"

// FileExists will verify if the file exists in the specified path.
func FileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}
