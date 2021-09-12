package sendo

type SMSOptions struct {
	TemplateName string
	Sender       string
	PhoneNumber  string
	Data         interface{}
	Extra        []interface{} // extra options for various implementations.
}

type VerificationOptions struct {
	TemplateName string
	Sender       string
	PhoneNumber  string
	Code         string        // in the send speedy we can send code.
	Extra        []interface{} // extra options for various implementations.
}

type SMSService interface {
	Send(o SMSOptions) error
	SendVerificationCode(o VerificationOptions) error
}
