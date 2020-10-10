package sendo

type EmailSender struct {
	Name  string // optional
	Email string // required
}

type EmailFrom struct {
	Name  string // optional
	Email string // required
}

type EmailTo struct {
	Name  string // optional
	Email string // required
}

type EmailCC struct {
	Name  string // optional
	Email string // required
}

type EmailReplyTo struct {
	Name  string // optional
	Email string // required
}

type SendSMTPEmailOptions struct {
	Sender       *EmailSender           // optional, driver should support default sender per for each template
	From         *EmailFrom             // for some drivers can be optional, for some other required.
	To           []EmailTo              // required
	CC           []EmailCC              // optional
	ReplyTo      *EmailReplyTo          // optional
	Subject      *string                // optional, drivers must support default subject for each template.
	TemplateName string                 // required
	Data         map[string]interface{} // optional
	Extra        []interface{}          // optional extra options.
	// TODO: add attachment,...
}

type EmailService interface {
	SendSMTP(o SendSMTPEmailOptions) error
}
