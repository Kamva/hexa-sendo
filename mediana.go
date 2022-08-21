package sendo

import (
	"bytes"
	"context"
	"text/template"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hexahttp"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
)

// medianaService implements the SMSService.
type medianaService struct {
	client        *hexahttp.Client
	defaultSender string
	token         string
	tpl           *template.Template
}

type MedianaOptions struct {
	ApiClient     *hexahttp.Client // APIClient must have base api url as its base API url. (e.g., abc.com/api/v1).
	Token         string
	DefaultSender string
	Templates     map[string]string
}

func NewMedianaService(o MedianaOptions) (SMSService, error) {
	t, err := parseTextTemplates("mediana_root", o.Templates)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	return &medianaService{
		client:        o.ApiClient,
		defaultSender: o.DefaultSender,
		token:         o.Token,
		tpl:           t,
	}, tracer.Trace(err)
}

type MedianaExtra struct {
	Token string
}

func (s medianaService) Send(_ context.Context, o SMSOptions) error {
	msg, err := s.renderTemplate(o.TemplateName, o.Data)
	if err != nil {
		return tracer.Trace(err)
	}

	token := s.token
	if extraOptToken := extraOption[*ExtraOptionToken](o.Extra); extraOptToken != nil {
		token = extraOptToken.Tokens
	}
	var authorizationHeader = hexahttp.AuthenticateHeader("apikey", "", token)

	sender := o.Sender
	if sender == "" {
		sender = s.defaultSender
	}

	payload := hexa.Map{
		"recipient": []string{"+" + o.PhoneNumber},
		"sender":    sender,
		"message":   msg,
	}

	_, err = hexahttp.ResponseErr(s.client.PostJsonForm("sms/send/webservice/single", payload, authorizationHeader))
	return tracer.Trace(err)
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
	hlog.Warn("mediana driver doesn't support verification code sms type.")
	return nil
}

var _ SMSService = &medianaService{}
