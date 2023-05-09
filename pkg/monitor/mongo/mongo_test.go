package mongo_test

import (
	"bcfmonitor/pkg/config"
	"bcfmonitor/pkg/monitor/mongo"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMongo(t *testing.T) {
	cfg, err := config.Load("../../../conf/dev.yaml")
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.NotEmpty(t, cfg.Databases, "No databases in configuration")

	for _, m := range cfg.Databases {
		t.Logf("Checking database: %s -> %s:%d", m.Name, m.Host, m.Port)
		s := mongo.NewService(m.Name, m.Host, m.Port, m.SSL)
		assert.NoError(t, s.Check(), "Error from database", m.Name)
	}
}