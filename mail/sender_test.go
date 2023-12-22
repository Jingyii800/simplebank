package mail

import (
	"testing"

	"github.com/Jingyii800/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	// skip this test in CI
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "Test Email"
	content := `
	<h1> Hello </h1>
	<p> Test message from Simple Bank server </p>
	`
	to := []string{"jiajj052@outlook.com"}

	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)

	require.NoError(t, err)
}
