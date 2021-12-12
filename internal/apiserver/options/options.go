package options

import (
	"encoding/json"

	genericoptions "github.com/rose839/IAM/internal/pkg/options"
	"github.com/rose839/IAM/internal/pkg/server"
	cliflag "github.com/rose839/IAM/pkg/app"
	"github.com/rose839/IAM/pkg/idutil"
)

// Options runs a iam api server.
type Options struct {
	GenericServerRunOptions *genericoptions.ServerRunOptions       `json:"server"   mapstructure:"server"`
	InsecureServing         *genericoptions.InsecureServingOptions `json:"insecure" mapstructure:"insecure"`
	SecureServing           *genericoptions.SecureServingOptions   `json:"secure"   mapstructure:"secure"`
	GRPCOptions             *genericoptions.GRPCOptions            `json:"grpc"     mapstructure:"grpc"`
	FeatureOptions          *genericoptions.FeatureOptions         `json:"feature"  mapstructure:"feature"`
	JwtOptions              *genericoptions.JwtOptions
	MySQLOptions            *genericoptions.MySQLOptions
}

// NewOptions creates a new Options object with default parameters.
func NewOptions() *Options {
	o := &Options{
		GenericServerRunOptions: genericoptions.NewServerRunOptions(),
		InsecureServing:         genericoptions.NewInsecureServingOptions(),
		SecureServing:           genericoptions.NewSecureServingOptions(),
		GRPCOptions:             genericoptions.NewGRPCOptions(),
		FeatureOptions:          genericoptions.NewFeatureOptions(),
		JwtOptions:              genericoptions.NewJwtOptions(),
		MySQLOptions:            genericoptions.NewMySQLOptions(),
	}

	return o
}

// ApplyTo applies the run options to the method receiver and returns self.
func (o *Options) ApplyTo(c *server.Config) error {
	return nil
}

// Flags returns flags for a specific APIServer by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.GenericServerRunOptions.AddFlags(fss.FlagSet("generic"))
	o.InsecureServing.AddFlags(fss.FlagSet("insecure serving"))
	o.SecureServing.AddFlags(fss.FlagSet("secure serving"))
	o.GRPCOptions.AddFlags(fss.FlagSet("grpc"))
	o.FeatureOptions.AddFlags(fss.FlagSet("features"))
	o.JwtOptions.AddFlags(fss.FlagSet("jwt"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))

	return fss
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)
	return string(data)
}

// Complete set default Options.
func (o *Options) Complete() error {
	if o.JwtOptions.Key == "" {
		o.JwtOptions.Key = idutil.NewSecretKey()
	}
	return o.SecureServing.Complete()
}

// Validate checks Options and return a slice of found errs.
func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.GenericServerRunOptions.Validate()...)
	errs = append(errs, o.InsecureServing.Validate()...)
	errs = append(errs, o.SecureServing.Validate()...)
	errs = append(errs, o.GRPCOptions.Validate()...)
	errs = append(errs, o.FeatureOptions.Validate()...)
	errs = append(errs, o.JwtOptions.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)

	return errs
}
