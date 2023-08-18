package liblog

import (
	"context"

	"github.com/sirupsen/logrus"
)

const cDefaultCorrelationID = "00000000000000000000"

func Default() *logrus.Entry {
	return logrus.WithField("x-correlation-id", cDefaultCorrelationID)
}

func WithContext(ctx context.Context) *logrus.Entry {
	return logrus.WithContext(ctx).WithField("x-correlation-id", getCorrelationID(ctx))
}
