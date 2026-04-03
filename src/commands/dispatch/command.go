package dispatch

import "context"

type Command interface {
	Path() []string
	Summary() string
	Run(ctx context.Context, args []string) error
}

type FuncCommand struct {
	commandPath []string
	commandHelp string
	run         func(ctx context.Context, args []string) error
}

func NewFuncCommand(path []string, summary string, run func(ctx context.Context, args []string) error) FuncCommand {
	return FuncCommand{commandPath: path, commandHelp: summary, run: run}
}

func (c FuncCommand) Path() []string {
	return append([]string{}, c.commandPath...)
}

func (c FuncCommand) Summary() string {
	return c.commandHelp
}

func (c FuncCommand) Run(ctx context.Context, args []string) error {
	return c.run(ctx, args)
}
