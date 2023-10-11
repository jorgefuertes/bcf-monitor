package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	cfg, err := Load("../../conf/dev.yaml")
	require.NoError(t, err)
	dump, err := cfg.Dump()
	require.NoError(t, err)
	t.Log(dump)
}
