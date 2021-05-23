package config

import (
	"errors"
	"os"
)

// global args
var (
	ServerAddress string
	IndexName     string
	QueryExp      string

	JoinExp string
	ExpFile string

	Verbose    bool
	OutputFile string
)

// search args
var (
	Size   int
	Offset int
)

func Validate() error {
	if ServerAddress == "" {
		return errors.New("server is required")
	}

	if IndexName == "" {
		return errors.New("idx is required")
	}

	if _, err := os.Stat(ExpFile); err != nil {
		return err
	}

	return nil
}
