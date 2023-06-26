package config

import (
	"bytes"
	"strconv"
)

const (
	DevEnvironment Environment = iota
	ProdEnvironment
)

type Environment uint

func (e *Environment) UnmarshalText(text []byte) error {
	src := string(bytes.ToLower(bytes.TrimSpace(text)))
	unquoted, err := strconv.Unquote(src)
	if err != nil {
		unquoted = src
	}

	switch unquoted {
	case "prod", "production":
		*e = ProdEnvironment
	default:
		*e = DevEnvironment
	}

	return nil
}
