package options

import (
	"encoding/json"

	genericoptions "github.com/rose839/IAM/internal/pkg/options"
	"github.com/rose839/IAM/internal/pkg/server"
	cliflag "github.com/rose839/IAM/pkg/app"
	"github.com/rose839/IAM/pkg/log"
)

// Options runs a authzserver.
type Options struct {
	RPCServer               string
	RootCA                  string
	GenericServerRunOptions *genericoptions.ServerRunOptions
	InsecureServing         *genericoptions.InsecureServingOptions
	SecureServing           *genericoptions.SecureServingOptions
	RedisOptions            *genericoptions.RedisOptions
	FeatureOptions          *genericoptions.FeatureOptions
	Log                     *log.Options
}

// NewOptions creates a new Options object with default parameters.
func NewOptions() *Options {
	o := Options{
		RPCServer:               "127.0.0.1:8081",
		RootCA:                  "",
		GenericServerRunOptions: genericoptions.NewServerRunOptions(),
		InsecureServing:         genericoptions.NewInsecureServingOptions(),
		SecureServing:           genericoptions.NewSecureServingOptions(),
		RedisOptions:            genericoptions.NewRedisOptions(),
		FeatureOptions:          genericoptions.NewFeatureOptions(),
		Log:                     log.NewOptions(),
	}

	return &o
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
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	o.FeatureOptions.AddFlags(fss.FlagSet("features"))
	o.Log.AddFlags(fss.FlagSet("logs"))

	fs := fss.FlagSet("misc")
	fs.StringVar(&o.RPCServer, "rpcserver", o.RPCServer, "The address of iam rpc server. "+
		"The rpc server can provide all the secrets and policies to use.")

	fs.StringVar(&o.RootCA, "client-ca-file", o.RootCA, ""+
		"If set, any request presenting a client certificate signed by one of "+
		"the authorities in the client-ca-file is authenticated with an identity "+
		"corresponding to the CommonName of the client certificate.")

	return fss
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}

// Complete set default Options.
func (o *Options) Complete() error {
	return o.SecureServing.Complete()
}

func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.GenericServerRunOptions.Validate()...)
	errs = append(errs, o.InsecureServing.Validate()...)
	errs = append(errs, o.SecureServing.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)
	errs = append(errs, o.FeatureOptions.Validate()...)
	errs = append(errs, o.Log.Validate()...)

	return errs
}
