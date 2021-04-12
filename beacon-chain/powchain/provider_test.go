package powchain

import (
	"testing"

	"github.com/prysmaticlabs/prysm/shared/httputils/authorizationmethod"
	"github.com/prysmaticlabs/prysm/shared/testutil/assert"
	logTest "github.com/sirupsen/logrus/hooks/test"
)

func TestHttpEndpoint(t *testing.T) {
	hook := logTest.NewGlobal()
	url := "http://test"

	t.Run("URL", func(t *testing.T) {
		endpoint := HttpEndpoint(url)
		assert.Equal(t, url, endpoint.Url)
		assert.Equal(t, authorizationmethod.None, endpoint.Auth.Method)
	})
	t.Run("URL with separator", func(t *testing.T) {
		endpoint := HttpEndpoint(url + ",")
		assert.Equal(t, url, endpoint.Url)
		assert.Equal(t, authorizationmethod.None, endpoint.Auth.Method)
	})
	t.Run("Basic auth", func(t *testing.T) {
		endpoint := HttpEndpoint(url + ",Basic username:password")
		assert.Equal(t, url, endpoint.Url)
		assert.Equal(t, authorizationmethod.Basic, endpoint.Auth.Method)
		assert.Equal(t, "dXNlcm5hbWU6cGFzc3dvcmQ=", endpoint.Auth.Value)
	})
	t.Run("Basic auth with incorrect format", func(t *testing.T) {
		hook.Reset()
		endpoint := HttpEndpoint(url + ",Basic username:password foo")
		assert.Equal(t, url, endpoint.Url)
		assert.Equal(t, authorizationmethod.None, endpoint.Auth.Method)
		assert.LogsContain(t, hook, "Skipping authorization")
	})
	t.Run("Bearer auth", func(t *testing.T) {
		endpoint := HttpEndpoint(url + ",Bearer token")
		assert.Equal(t, url, endpoint.Url)
		assert.Equal(t, authorizationmethod.Bearer, endpoint.Auth.Method)
		assert.Equal(t, "token", endpoint.Auth.Value)
	})
	t.Run("Bearer auth with incorrect format", func(t *testing.T) {
		hook.Reset()
		endpoint := HttpEndpoint(url + ",Bearer token foo")
		assert.Equal(t, url, endpoint.Url)
		assert.Equal(t, authorizationmethod.None, endpoint.Auth.Method)
		assert.LogsContain(t, hook, "Skipping authorization")
	})
	t.Run("Too many separators", func(t *testing.T) {
		endpoint := HttpEndpoint(url + ",Bearer token,foo")
		assert.Equal(t, url, endpoint.Url)
		assert.Equal(t, authorizationmethod.None, endpoint.Auth.Method)
		assert.LogsContain(t, hook, "Skipping authorization")
	})
}
