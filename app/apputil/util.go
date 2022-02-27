package apputil

import "os"

// SetEnvs to os
func SetEnvs(mp map[string]string) {
	for key, value := range mp {
		_ = os.Setenv(key, value)
	}
}
