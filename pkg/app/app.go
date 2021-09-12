package app

// App is the main structure of a cli app, it's recommanded
// that an app be created with the app.NewApp func.
type App struct {
	basename    string
	name        string
	description string
	runFunc     RunFunc
	slience     bool
	noVersion   bool
	noConfig    bool
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

// WithSilence sets the application to silent mode, in which the program startup
// information, configuration information, and version information are not
// printed in the console.
func WithSilence() Option {
	return func(app *App) {
		app.slience = true
	}
}

// WithNoVersion set the application does not provide version flag.
func WithNoVersion() Option {
	return func(a *App) {
		a.noVersion = true
	}
}

// WithNoConfig set the application does not provide config flag.
func WithNoConfig() Option {
	return func(a *App) {
		a.noConfig = true
	}
}

func NewApp(name string, basename string, opts ...Option) *App {
	a := &App{
		name:     name,
		basename: basename,
	}

	for _, o := range opts {
		o(a)
	}

	return a
}
