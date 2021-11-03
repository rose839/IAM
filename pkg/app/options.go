package app

import "github.com/spf13/pflag"

// CliOptions abstracts app-defined options for reading parameters from the
// command line.
type CliOptions interface {
	AddFlags(fs *pflag.FlagSet) // AddFlags adds flags to the specified FlagSet object.
	Flags() NamedFlagSets       // get the specified name flagset
	Validate() []error          // check whether option args is validate
}

// ConfigurableOptions abstracts app-defined options for reading parameters
// from a configuration file.
type ConfigurableOptions interface {
	// ApplyFlags parsing parameters from the command line or configuration file
	// to the app-defined options struct instance.
	ApplyFlags() []error
}

// CompleteableOptions abstracts app-defined options which can be completed.
type CompleteableOptions interface {
	// If option args could be completed, the complete it.
	// Otherwise return error
	Complete() error
}

// PrintableOptions abstracts app-defined options which can be printed.
type PrintableOptions interface {
	String() string
}
