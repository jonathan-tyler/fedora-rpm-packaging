package platform

import "context"

type Runner interface {
	Run(ctx context.Context, command Command) error
	Capture(ctx context.Context, command Command) (string, error)
}
