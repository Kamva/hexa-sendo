package sendo

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hexahttp"
	"github.com/kamva/tracer"
	"text/template"
)

// medianaUserAppService implements the SMSService.
type medianaUserAppService struct {
	client *hexahttp.Client
	tpl    *template.Template
}

type MedianaUserAppOptions struct {
	APIUrl        string
	Token         string
	DefaultSender string
	Templates     map[string]string
}

func NewMedianaUserAppService(o MedianaUserAppOptions) (SMSService, error) {
	t, err := parseTextTemplates("mediana_root", o.Templates)

	return &medianaService{
		client: hexahttp.NewClient(&o.APIUrl),
		tpl:    t,
	}, tracer.Trace(err)
}

type MedianaExtra struct {
	Token string
}

func (s medianaUserAppService) Send(_ context.Context, o SMSOptions) error {
	msg, err := s.renderTemplate(o.TemplateName, o.Data)
	if err != nil {
		return tracer.Trace(err)
	}

	extraJsonStr, err := json.Marshal(o.Extra)
	if err != nil {
		return tracer.Trace(err)
	}
	var medianaExtra MedianaExtra
	if err = json.Unmarshal(extraJsonStr, &medianaExtra); err != nil {
		return tracer.Trace(err)
	}

	var authorizationHeader = hexahttp.AuthenticateHeader("apikey", "", medianaExtra.Token)
	var recipients []string
	recipients = append(recipients, "+"+o.PhoneNumber)
	resp, err := s.client.PostJsonFormWithOptions("sms/send/webservice/single", hexa.Map{
		"recipient": recipients,
		"sender":    o.Sender,
		"message":   msg,
	}, authorizationHeader)
	if err != nil {
		return err
	}
	return hexahttp.ResponseError(resp)
}

func (s *medianaUserAppService) renderTemplate(tplName string, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := s.tpl.ExecuteTemplate(&buf, tplName, data); err != nil {
		return "", tracer.Trace(err)
	}
	return buf.String(), nil
}

// SendVerificationCode ignores the sender param in mediana driver.
func (s medianaUserAppService) SendVerificationCode(_a context.Context, o VerificationOptions) error {

	return nil
}

var _ SMSService = &medianaUserAppService{}
