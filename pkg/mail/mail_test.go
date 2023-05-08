package mail_test

import (
	"bcfmonitor/pkg/config"
	"bcfmonitor/pkg/mail"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMail(t *testing.T) {
	cfg, err := config.Load("../../conf/dev.yaml")
	require.NoError(t, err)
	require.NotNil(t, cfg)
	ms := mail.NewService(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.User, cfg.SMTP.Pass)
	err = ms.Send("jorge@jorgefuertes.com", "Prueba de correo.", "Correo de test, borrar.")
	require.NoError(t, err)
}
