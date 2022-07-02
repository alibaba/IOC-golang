package config

import (
	"regexp"
)

func isEnv(envValue string) bool {
	// ${ Xxx_Yyy_Zzz }
	ok, err := regexp.Match("^\\$\\{[A-Z_]+}$", []byte(envValue))
	if err != nil || !ok {
		return false
	}

	return true
}
