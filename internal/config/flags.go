package config

import "github.com/spf13/pflag"

type CommandFlags struct {
	EnvFile string
}

func ReadCommandFlags() CommandFlags {
	var params CommandFlags
	pflag.StringVarP(&params.EnvFile, "env-file", "f", "", "Environment file to load (.env)")
	pflag.Parse()
	return params
}
