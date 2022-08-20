package sendo

import (
	"bytes"
	"context"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hexahttp"
	"github.com/kamva/tracer"
	"text/template"
)

// medianaService implements the SMSService.
type medianaService struct {
	client        *hexahttp.Client
	defaultSender string
	token         string
	tpl           *template.Template
}

type MedianaOptions struct {
	APIUrl        string
	Token         string
	DefaultSender string
	Templates     map[string]string
}

func NewMedianaService(o MedianaOptions) (SMSService, error) {
	t, err := parseTextTemplates("kavenegar_root", o.Templates)

	return &medianaService{
		client:        hexahttp.NewClient(&o.APIUrl),
		defaultSender: o.DefaultSender,
		token:         o.Token,
		tpl:           t,
	}, tracer.Trace(err)
}

type medianaMessage struct {
	Message string
}

func (s medianaService) Send(_ context.Context, o SMSOptions) error {
	msg, err := s.renderTemplate(o.TemplateName, o.Data)
	if err != nil {
		return tracer.Trace(err)
	}
	var authorizationHeader = hexahttp.AuthenticateHeader("apikey", "", s.token)
	var recipients []string
	recipients = append(recipients, "+"+o.PhoneNumber)
	sender := o.Sender
	if sender == "" {
		sender = s.defaultSender
	}

	resp, err := s.client.PostJsonFormWithOptions("sms/send/webservice/single", hexa.Map{
		"recipient": recipients,
		"sender":    sender,
		"message":   msg,
	}, authorizationHeader)
	if err != nil {
		return err
	}
	return hexahttp.ResponseError(resp)
}

func (s *medianaService) renderTemplate(tplName string, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := s.tpl.ExecuteTemplate(&buf, tplName, data); err != nil {
		return "", tracer.Trace(err)
	}
	return buf.String(), nil
}

// SendVerificationCode ignores the sender param in mediana driver.
func (s medianaService) SendVerificationCode(_a context.Context, o VerificationOptions) error {

	return nil
}

var _ SMSService = &kavenegarService{}
