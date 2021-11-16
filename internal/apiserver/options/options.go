package options

import (
	"encoding/json"

	cliflag "github.com/rose839/pkg/app/flag"
)

// Options runs a iam api server.
type Options struct {
}

// NewOptions creates a new Options object with default parameters.
func NewOptions() *Options {
	o := &Options{}

	return o
}

// ApplyTo applies the run options to the method receiver and returns self.
func (o *Options) ApplyTo(server) error {
	return nil
}

// Flags returns flags for a specific APIServer by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {

}

func (o *Options) String() string {
	data, _ := json.Marshal(o)
	return string(data)
}

// Complete set default Options.
func (o *Options) Complete() error {

}

// Validate checks Options and return a slice of found errs.
func (o *Options) Validate() []error {
	var errs []error

	return errs
}
