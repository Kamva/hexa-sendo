package sendo

type SendSMSOptions struct {
	TemplateName string
	PhoneNumber  string
	Data         interface{}
	Extra        []interface{} // extra options for various implementations.
}

type SendSpeedySMSOptions struct {
	TemplateName string
	PhoneNumber  string
	Code         string        // in the send speedy we can send code.
	Extra        []interface{} // extra options for various implementations.
}

type SMSService interface {
	Send(o SendSMSOptions) error
	SendSpeedySMS(o SendSpeedySMSOptions) error
}
