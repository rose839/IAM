package config

import "github.com/rose839/IAM/internal/apiserver/options"

// Config is the running configuration structure of the IAM api service.
type Config struct {
	*options.Options
}

// CreateConfigFromOptions creates a running configuration instance based
// on a given IAM api command line or configuration file option.
func CreateConfigFromOptions(opts *options.Options) (*Config, error) {
	return &Config{opts}, nil
}
