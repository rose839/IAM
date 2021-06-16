package app

// App is the main structure of a cli app, it's recommanded
// that an app be created with the app.NewApp func.
type App struct {
	basename string
	name string
	description string
	runFunc RunFunc
}

// Option defines optional parameters for initializing the application
// structure.
type Option func(*App)

// RunFunc defines the application's startup callback function.
type RunFunc func(basename string) error

// WithRunFunc is used to set the application startup callback function option.
func WithRunFunc(run RunFunc) Option {
	return func(app *App) {
		app.runFunc = run
	}
}

// WithDescription is used to set the description of the application.
func WithDescription(desc string) Option {
	return func(app *App) {
		app.description = desc
	}
}

func NewApp(name string, basename string, opts... Option) *App {
	a := &App {
		name: name,
		basename: basename,
	}

	for _, o := range opts {
		o(a)
	}

	return a
}