package health

import "context"

type Checker interface {
	HealthCheck(ctx context.Context) error
}

type BootChecker interface {
	BootHealthCheck(ctx context.Context) error
}
