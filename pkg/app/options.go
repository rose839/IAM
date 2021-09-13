package app

import "github.com/spf13/pflag"

// CliOptions abstracts configuration options for reading parameters from the
// command line.
type CliOptions interface {
	AddFlags(fs *pflag.FlagSet) // AddFlags adds flags to the specified FlagSet object.
	Flags() NamedFlagSets       // get the specified name flagset
	Validate() []error
}

// ConfigurableOptions abstracts configuration options for reading parameters
// from a configuration file.
type ConfigurableOptions interface {
	// ApplyFlags parsing parameters from the command line or configuration file
	// to the options instance.
	ApplyFlags() []error
}

// CompleteableOptions abstracts options which can be completed.
type CompleteableOptions interface {
	Complete() error
}

// PrintableOptions abstracts options which can be printed.
type PrintableOptions interface {
	String() string
}
