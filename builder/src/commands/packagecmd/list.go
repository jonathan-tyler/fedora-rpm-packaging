package packagecmd

import (
	"context"
	"fmt"
	"io"

	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/app/bootstrap"
	"github.com/jonathan-tyler/fedora-rpm-packaging/builder/src/commands/dispatch"
)

func NewListCommand(services *bootstrap.Services, stdout io.Writer) dispatch.FuncCommand {
	return dispatch.NewFuncCommand([]string{"package", "list"}, "list configured packages", func(ctx context.Context, args []string) error {
		if len(args) != 0 {
			return fmt.Errorf("usage: fpb package list")
		}
		for _, definition := range services.Registry.List() {
			if _, err := fmt.Fprintln(stdout, definition.Name); err != nil {
				return err
			}
		}
		return nil
	})
}
