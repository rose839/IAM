// Package apiserver does all of the work necessary to create a iam APIServer.
package apiserver

import (
	"github.com/rose839/IAM/internal/apiserver/config"
	"github.com/rose839/IAM/internal/apiserver/options"
	"github.com/rose839/IAM/pkg/app"
)

const commandDesc = `The IAM API server validates and configures data
for the api objects which include users, policies, secrets, and
others. The API Server services REST operations to do the api objects management.`

func NewApp(basename string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp(
		"IAM API Server",
		basename,
		app.WithOptions(opts),
		app.WithDescription(commandDesc),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
	)

	return application
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		// init log

		// create app config
		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}

		return Run(cfg)
	}
}
