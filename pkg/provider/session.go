package provider

import "github.com/aws/aws-sdk-go/aws/session"

func NewSession() *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
}
