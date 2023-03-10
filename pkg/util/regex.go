package util

import (
	"regexp"
)

func RegretAllVariableName(s string) ([]string, error) {
	regexpCompile, err := regexp.Compile(`[a-zA-Z_][a-zA-Z0-9_]*`)
	if err != nil {
		return nil, err
	}
	match := regexpCompile.FindAllString(s, -1)
	return match, nil
}
