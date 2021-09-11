package sendo

import (
	"bytes"
	"github.com/kamva/tracer"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"text/template"
)

type TwilioSender string

type TwilioService struct {
	client        *twilio.RestClient
	tpl           *template.Template
	defaultSender TwilioSender
}

type TwilioOptions struct {
	AccountSID string
	AuthToken  string
	Templates  map[string]string
}

func NewTwilioService(o TwilioOptions) (SMSService, error) {

	t, err := parseTextTemplates("twilio_root", o.Templates)

	return &TwilioService{
		client: twilio.NewRestClientWithParams(twilio.RestClientParams{
			Username: o.AccountSID,
			Password: o.AuthToken,
		}),
		tpl: t,
	}, tracer.Trace(err)
}

func (s *TwilioService) renderTemplate(tplName string, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := s.tpl.ExecuteTemplate(&buf, tplName, data); err != nil {
		return "", tracer.Trace(err)
	}
	return buf.String(), nil
}

func (s *TwilioService) getSender(extraOptions []interface{}) TwilioSender {
	for _, v := range extraOptions {
		if s, ok := v.(TwilioSender); ok {
			return s
		}
	}
	return s.defaultSender
}

func (s *TwilioService) Send(o SMSOptions) error {
	msg, err := s.renderTemplate(o.TemplateName, o.Data)
	if err != nil {
		return tracer.Trace(err)
	}

	params := &openapi.CreateMessageParams{}
	params.SetTo(o.PhoneNumber)
	params.SetFrom(string(s.getSender(o.Extra)))
	params.SetBody(msg)

	_, err = s.client.ApiV2010.CreateMessage(params)
	return tracer.Trace(err)
}

func (s *TwilioService) SendVerificationCode(o VerificationOptions) error {
	msg, err := s.renderTemplate(o.TemplateName, map[string]interface{}{
		"code": o.Code,
	})

	if err != nil {
		return tracer.Trace(err)
	}

	params := &openapi.CreateMessageParams{}
	params.SetTo(o.PhoneNumber)
	params.SetFrom(string(s.getSender(o.Extra)))
	params.SetBody(msg)

	_, err = s.client.ApiV2010.CreateMessage(params)
	return tracer.Trace(err)
}

var _ SMSService = &TwilioService{}
