package idetool

import (
	"os"
)

// LoadHceFile file contents
func LoadHceFile(hceFile string) (bts []byte, err error) {
	return os.ReadFile(hceFile)
}
