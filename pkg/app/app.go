package app

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/rose839/IAM/pkg/log"
	"github.com/rose839/IAM/pkg/term"
	"github.com/rose839/IAM/pkg/version/verflag"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var progressMessage = color.GreenString("==>")

// App is the main structure of a cli app, it's recommanded
// that an app be created with the app.NewApp func.
type App struct {
	basename    string               // app name
	name        string               // short description
	description string               // long description
	runFunc     RunFunc              // user-defined main func
	silence     bool                 // silence mode at startup phase
	noConfig    bool                 // whether add "--config" flag
	noVersion   bool                 // whether add "--version" flag
	options     CliOptions           // app-defined struct of cli options, used to store config and register flag
	commands    []*Command           // sub commands
	args        cobra.PositionalArgs // validation func of positional arguments(arg is belong to command, not flag)
	cmd         *cobra.Command       // root command
}

// Option defines optional parameters for initializing the application structure.
type Option func(*App)

// WithOptions to open the application's function to read from the command line
// or read parameters from the configuration file for every app.
// Application shoule provide an user-defined structure that implements CliOptions to
// give cli flag and receive cli params.
func WithOptions(opt CliOptions) Option {
	return func(a *App) {
		a.options = opt
	}
}

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
		app.silence = true
	}
}

// WithNoConfig set the application does not provide "--config" flag.
func WithNoConfig() Option {
	return func(a *App) {
		a.noConfig = true
	}
}

// WithNoVersion set the application does not provide "--version" flag.
func WithNoVersion() Option {
	return func(a *App) {
		a.noVersion = true
	}
}

// WithValidArgs set the validation function to valid root command arguments.
func WithValidArgs(args cobra.PositionalArgs) Option {
	return func(a *App) {
		a.args = args
	}
}

// WithDefaultValidArgs set default validation function to valid valid non-flag arguments.
func WithDefaultValidArgs() Option {
	return func(a *App) {
		a.args = func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		}
	}
}

// WithCommands add sub-commands to root command.
func WithCommands(commands ...*Command) Option {
	return func(a *App) {
		a.commands = commands
	}
}

// NewApp creates a new application instance based on the given application name,
// base name, and other options.
func NewApp(name string, basename string, opts ...Option) *App {
	a := &App{
		name:     name,
		basename: basename,
	}

	// Set options
	for _, o := range opts {
		o(a)
	}

	a.buildCommand()

	return a
}

// buildCommand is used to set cobra command and flags.
func (a *App) buildCommand() {
	// Create root command
	cmd := cobra.Command{
		Use:           FormatBaseName(a.basename), // Add basename(app name) to one-line usage
		Short:         a.name,
		Long:          a.description,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          a.args,
	}

	cmd.SetOut(os.Stdout)        // Set destination for usage message
	cmd.SetErr(os.Stderr)        // Set destination for cmd error message
	cmd.Flags().SortFlags = true // Sort flag by flag name
	initFlag(cmd.Flags())

	// Add sub command to root command
	if len(a.commands) > 0 {
		for _, command := range a.commands {
			cmd.AddCommand(command.CobraCommand())
		}

		// Add "help" sub command
		cmd.SetHelpCommand(helpCommand(a.name))
	}

	// Add flags to root command
	var namedFlagSets NamedFlagSets
	if a.options != nil {
		// set app-defined grouped flagset
		namedFlagSets = a.options.Flags()
		fs := cmd.Flags()
		for _, f := range namedFlagSets.FlagSets {
			fs.AddFlagSet(f) // add flag
		}

		usageFmt := "Usage:\n  %s\n"
		termWidth, _, err := term.TerminalSize(cmd.OutOrStdout())
		if err != nil {
			log.Fatalf(err.Error())
		}

		// set help func, called at "-h/--help" flag
		cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
			PrintSections(cmd.OutOrStdout(), namedFlagSets, termWidth)
		})

		// set usage func, called at error input,
		// could shut down output usage by set SilenceUsage: true
		cmd.SetUsageFunc(func(cmd *cobra.Command) error {
			fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
			PrintSections(cmd.OutOrStderr(), namedFlagSets, termWidth)

			return nil
		})
	}

	// Add "--config" flag
	if !a.noConfig {
		addConfigFlag(a.basename, namedFlagSets.FlagSet("global"))
	}

	// Add "--version" flag
	if !a.noVersion {
		verflag.AddFlags(namedFlagSets.FlagSet("global"))
	}

	// Add "-h/--help" flag
	namedFlagSets.FlagSet("global").BoolP("help", "h", false, fmt.Sprintf("help for %s", cmd.Name()))

	// add new global flagset to cmd FlagSet
	cmd.Flags().AddFlagSet(namedFlagSets.FlagSet("global"))

	// set main app run func
	if a.runFunc != nil {
		cmd.RunE = a.runCommand
	}

	a.cmd = &cmd
}

// Run is used to launch the application.
func (a *App) Run() {
	if err := a.cmd.Execute(); err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error:"), err)
		os.Exit(1)
	}
}

// Command returns cobra command instance inside the application.
func (a *App) Command() *cobra.Command {
	return a.cmd
}

// Command main func, will call user-defined run func.
func (a *App) runCommand(cmd *cobra.Command, args []string) error {
	printWorkingDir()
	PrintFlag(cmd.Flags())

	if !a.noVersion {
		// display application version information
		verflag.PrintAndExitIfRequested()
	}

	if !a.noConfig {
		// Bind viper config with command flags
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}

		// Store viper config to app-defined CliOptions struct
		if err := viper.Unmarshal(a.options); err != nil {
			return err
		}
	}

	// Print some info
	if !a.silence {
		log.Infof("%v Starting %s ...", progressMessage, a.name)
		if !a.noVersion {

		}
		if !a.noConfig {
			log.Infof("%v Config file used: `%s`", progressMessage, viper.ConfigFileUsed())
		}
	}

	// Complete/validate/print config options
	if a.options != nil {
		if err := a.applyOptionRules(); err != nil {
			return err
		}
	}

	// Call user-defined main func
	if a.runFunc != nil {
		return a.runFunc(a.basename)
	}

	return nil
}

func (a *App) applyOptionRules() error {
	// Complete config options
	if completeableOptions, ok := a.options.(CompleteableOptions); ok {
		if err := completeableOptions.Complete(); err != nil {
			return err
		}
	}

	// Validate config options
	if errs := a.options.Validate(); len(errs) != 0 {
		return errs[0]
	}

	// Print config options
	if printableOptions, ok := a.options.(PrintableOptions); ok && !a.silence {
		log.Infof("%v Config: `%s`", progressMessage, printableOptions.String())
	}

	return nil
}

func printWorkingDir() {
	wd, _ := os.Getwd()
	log.Infof("%v WorkingDir: %s", progressMessage, wd)
}
