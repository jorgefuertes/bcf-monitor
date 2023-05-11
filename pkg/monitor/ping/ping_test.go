package ping_test

import (
	"bcfmonitor/pkg/config"
	"bcfmonitor/pkg/monitor/ping"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedis(t *testing.T) {
	cfg, err := config.Load("../../../conf/dev.yaml")
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.NotEmpty(t, cfg.Pings, "No pings in configuration")

	for _, p := range cfg.Pings {
		t.Logf("Checking web: %s", p.Name)
		s := ping.NewService(p.Name, p.Host, p.Timeout, p.Every)
		if s.Name() == "unreachable" {
			assert.Error(t, s.Check())
			continue
		}
		assert.NoError(t, s.Check(), "Error from ping %s", p.Name)
	}
}
