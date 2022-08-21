package sendo

import (
	"bytes"
	"context"
	"text/template"

	"github.com/kamva/tracer"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioService struct {
	client        *twilio.RestClient
	tpl           *template.Template
	defaultSender string
}

type TwilioOptions struct {
	AccountSID    string
	AuthToken     string
	Templates     map[string]string
	DefaultSender string
}

func NewTwilioService(o TwilioOptions) (SMSService, error) {
	t, err := parseTextTemplates("twilio_root", o.Templates)

	return &TwilioService{
		tpl: t,
		client: twilio.NewRestClientWithParams(twilio.RestClientParams{
			Username: o.AccountSID,
			Password: o.AuthToken,
		}),
		defaultSender: o.DefaultSender,
	}, tracer.Trace(err)
}

func (s *TwilioService) renderTemplate(tplName string, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := s.tpl.ExecuteTemplate(&buf, tplName, data); err != nil {
		return "", tracer.Trace(err)
	}
	return buf.String(), nil
}

func (s *TwilioService) Send(_ context.Context, o SMSOptions) error {
	msg, err := s.renderTemplate(o.TemplateName, o.Data)
	if err != nil {
		return tracer.Trace(err)
	}

	if o.Sender == "" {
		o.Sender = s.defaultSender
	}

	params := &openapi.CreateMessageParams{}
	params.SetTo(o.PhoneNumber)
	params.SetFrom(o.Sender)
	params.SetBody(msg)

	_, err = s.client.ApiV2010.CreateMessage(params)
	return tracer.Trace(err)
}

func (s *TwilioService) SendVerificationCode(_ context.Context, o VerificationOptions) error {
	msg, err := s.renderTemplate(o.TemplateName, map[string]interface{}{
		"code": o.Code,
	})

	if err != nil {
		return tracer.Trace(err)
	}

	if o.Sender == "" {
		o.Sender = s.defaultSender
	}

	params := &openapi.CreateMessageParams{}
	params.SetTo(o.PhoneNumber)
	params.SetFrom(o.Sender)
	params.SetBody(msg)

	_, err = s.client.ApiV2010.CreateMessage(params)
	return tracer.Trace(err)
}

var _ SMSService = &TwilioService{}
