package ast

import (
	"errors"
	"strings"
)

func parseWord(words []string) (string, int, error) {
	// TODO: parse words start with quate
	if len(words) == 0 {
		return "", 0, errors.New("no word to parse")
	}
	return strings.TrimSpace(words[0]), 1, nil
}

func isOperation(target string) bool {
	switch target {
	case termOper, termsOper, rangeOper:
		return true
	default:
		return false
	}
}
