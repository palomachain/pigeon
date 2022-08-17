package health

import "context"

type Checker interface {
	HealthCheck(ctx context.Context) error
}
