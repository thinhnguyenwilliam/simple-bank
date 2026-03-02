// simple-bank/mail/sender_test.go
// go test -v -count=1 ./mail
package mail

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thinhcompany/simple-bank/util"
)

func TestSendEmailWithGmail(t *testing.T) {
	// t.Skip("skip real email sending")

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(
		config.EmailSenderName,
		config.EmailSenderAddress,
		config.EmailSenderPassword,
	)

	err = sender.SendEmail(
		"Test Email 123",
		"<h1>Hello from Simple Bank 👋</h1>",
		[]string{"hoangproee@gmail.com"},
		nil,
		nil,
		nil,
	)

	require.NoError(t, err)
}
