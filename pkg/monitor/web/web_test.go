package web_test

import (
	"bcfmonitor/pkg/config"
	"bcfmonitor/pkg/monitor/web"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedis(t *testing.T) {
	cfg, err := config.Load("../../../conf/dev.yaml")
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.NotEmpty(t, cfg.Caches, "No caches in configuration")

	for _, w := range cfg.Webs {
		t.Logf("Checking web: %s", w.Name)
		s := web.NewService(w.Name, w.URL, w.Needle, w.HeaderMap(), w.Timeout, w.Every)
		assert.NoError(t, s.Check(), "Error from web", w.Name)
	}
}