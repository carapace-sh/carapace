package config

import "os"

func IsLenient() (lenient bool) {
	_, lenient = os.LookupEnv("CARAPACE_LENIENT")
	return
}
