package blxr

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func (c *Client) KeepAlive(ctx context.Context, _ sync.Locker) error {
	return c.runHealthCheck(ctx)
}

func (c *Client) runHealthCheck(ctx context.Context) error {
	res, err := c.rs.R().SetContext(ctx).
		SetHeader("Authorization", c.authHeader).
		SetBody(map[string]interface{}{"method": "ping", "id": "1", "params": nil}).
		Post(cBloXRouteCloudAPIURL)

	if err != nil || res.StatusCode() != http.StatusOK {
		if c.IsHealthy() {
			log.WithContext(ctx).
				WithField("status code", res.StatusCode()).
				WithError(err).
				Warnf("Blxr client lost connection. Pigeon EVM relayer traits unhealthy.")
			c.isHealthy = false
		}
		return fmt.Errorf("BLXR client unhealthy: %w", err)
	}

	if !c.IsHealthy() {
		log.WithContext(ctx).Infof("Blxr client recovered.")
		c.isHealthy = true
	}

	return nil
}

func (c *Client) GetHealthprobeInterval() time.Duration {
	return cHealthprobeQueryInterval
}
