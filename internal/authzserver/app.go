package authzserver

import (
	"github.com/rose839/IAM/internal/authzserver/config"
	"github.com/rose839/IAM/internal/authzserver/options"
	"github.com/rose839/IAM/pkg/app"
	"github.com/rose839/IAM/pkg/log"
)

const commandDesc = `Authorization server to run ladon policies which can protecting your resources.`

func NewApp(basename string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp("IAM Authorization Server",
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
		log.Init(opts.Log)
		defer log.Flush()

		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}

		return Run(cfg)
	}
}
