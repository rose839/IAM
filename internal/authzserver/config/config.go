package config

import "github.com/rose839/IAM/internal/apiserver/options"

// Config is the running configuration structure of the IAM authzserver service.
type Config struct {
	*options.Options
}

func CreateConfigFromOptions(opts *options.Options) (*Config, error) {
	return &Config{opts}, nil
}
