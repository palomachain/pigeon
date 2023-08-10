package blxr

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestHealthProbe(t *testing.T) {
	expectedBody := `{"id":"1","method":"ping","params":null}`
	fixture := `{"id":"1","result":{"pong":"2023-08-09 15:29:02.467176"},"jsonrpc":"2.0"}`
	responder := httpmock.NewStringResponder(200, fixture)
	httpmock.RegisterMatcherResponder("POST", cBloXRouteCloudAPIURL, httpmock.BodyContainsString(expectedBody), responder)
	authHeather := "0xDEADBEEF"

	c := New(authHeather)
	httpmock.ActivateNonDefault(c.rs.GetClient())

	err := c.runHealthCheck(context.Background())
	require.NoError(t, err)

	count := httpmock.GetTotalCallCount()

	require.Equal(t, 1, count)
}
