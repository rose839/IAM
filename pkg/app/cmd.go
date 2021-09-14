package app

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Command is a sub command structure of a cli application.
// It is recommended that a command be created with the app.NewCommand()
// function.
type Command struct {
	usage    string
	desc     string
	options  CliOptions
	commands []*Command     // nested sub command
	runFunc  RunCommandFunc // user-defined sub-command main func
}

// CommandOption defines optional parameters for initializing the command
// structure.
type CommandOption func(*Command)

// WithCommandOptions to open the application's function to read from the
// command line.
func WithCommandOptions(opt CliOptions) CommandOption {
	return func(c *Command) {
		c.options = opt
	}
}

// RunCommandFunc defines the application's command startup callback function.
type RunCommandFunc func(args []string) error

// WithCommandRunFunc is used to set the application's command startup callback
// function option.
func WithCommandRunFunc(run RunCommandFunc) CommandOption {
	return func(c *Command) {
		c.runFunc = run
	}
}

// NewCommand creates a new sub command instance based on the given command name
// and other options.
func NewCommand(usage string, desc string, opts ...CommandOption) *Command {
	c := &Command{
		usage: usage,
		desc:  desc,
	}

	for _, o := range opts {
		o(c)
	}

	return c
}

// AddCommand adds sub command to the current command.
func (c *Command) AddCommand(cmd *Command) {
	c.commands = append(c.commands, cmd)
}

// AddCommands adds multiple sub commands to the current command.
func (c *Command) AddCommands(cmds ...*Command) {
	c.commands = append(c.commands, cmds...)
}

func (c *Command) CobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   c.usage,
		Short: c.desc,
	}

	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.Flags().SortFlags = true
	initFlag(cmd.Flags())

	// Add sub commands
	if len(c.commands) > 0 {
		for _, command := range c.commands {
			cmd.AddCommand(command.CobraCommand())
		}
	}

	if c.runFunc != nil {
		cmd.Run = c.runCommand
	}

	// Add flagset to command
	if c.options != nil {
		for _, f := range c.options.Flags().FlagSets {
			cmd.Flags().AddFlagSet(f)
		}
	}

	addHelpCommandFlag(c.usage, cmd.Flags())

	return cmd
}

func (c *Command) runCommand(cmd *cobra.Command, args []string) {
	if c.runFunc != nil {
		if err := c.runFunc(args); err != nil {
			fmt.Printf("%v %v\n", color.RedString("Error:"), err)
			os.Exit(1)
		}
	}
}

// AddCommand adds sub command to the application.
func (a *App) AddCommand(cmd *Command) {
	a.commands = append(a.commands, cmd)
}

// AddCommands adds multiple sub commands to the application.
func (a *App) AddCommands(cmds ...*Command) {
	a.commands = append(a.commands, cmds...)
}