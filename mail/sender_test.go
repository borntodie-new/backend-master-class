package mail

import (
	"github.com/borntodie-new/backend-master-class/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSenderEmailWithQQ(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewQQSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	subject := "A test email"
	content := `
	<h1>Hello world</h1>
	<p>This is a test message from <strong>Master Backend Class</strong></p>
	`
	to := []string{"YourName@example.com"}
	attachFiles := []string{"../README.md"}
	err = sender.SendEmail(subject, content, to, attachFiles)
	require.NoError(t, err)
}
