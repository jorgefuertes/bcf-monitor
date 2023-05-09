package redis_test

import (
	"bcfmonitor/pkg/config"
	"bcfmonitor/pkg/monitor/redis"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedis(t *testing.T) {
	cfg, err := config.Load("../../../conf/dev.yaml")
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.NotEmpty(t, cfg.Caches, "No caches in configuration")

	for _, r := range cfg.Caches {
		t.Logf("Checking cache: %s -> %s:%d", r.Name, r.Host, r.Port)
		s := redis.NewService(r.Name, r.Host, r.Port, r.Password)
		assert.NoError(t, s.Check(), "Error from cache", r.Name)
	}
}