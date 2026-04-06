package dispatch

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
)

type Dispatcher struct {
	commands []Command
}

func New(commands ...Command) Dispatcher {
	return Dispatcher{commands: commands}
}

func (d Dispatcher) Execute(ctx context.Context, args []string) error {
	if len(args) == 0 || isHelp(args) {
		return errors.New(d.Usage())
	}

	for _, command := range d.commands {
		path := command.Path()
		if hasPrefix(args, path) {
			return command.Run(ctx, args[len(path):])
		}
	}

	return fmt.Errorf("unknown command: %s\n\n%s", strings.Join(args, " "), d.Usage())
}

func (d Dispatcher) Usage() string {
	rows := make([]string, 0, len(d.commands))
	for _, command := range d.commands {
		rows = append(rows, fmt.Sprintf("  %-28s %s", strings.Join(command.Path(), " "), command.Summary()))
	}
	sort.Strings(rows)

	return "Usage: fpb <group> <verb> [args]\n\nCommands:\n" + strings.Join(rows, "\n") + "\n"
}

func hasPrefix(args []string, prefix []string) bool {
	if len(args) < len(prefix) {
		return false
	}
	for index := range prefix {
		if args[index] != prefix[index] {
			return false
		}
	}
	return true
}

func isHelp(args []string) bool {
	return len(args) == 1 && (args[0] == "help" || args[0] == "--help" || args[0] == "-h")
}
