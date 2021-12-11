package sendo

import (
	"bytes"
	"context"
	"text/template"

	"github.com/kamva/tracer"
	"github.com/kavenegar/kavenegar-go"
)

// kavenegarService implements the SMSService.
type kavenegarService struct {
	client        *kavenegar.Kavenegar
	defaultSender string

	tpl *template.Template
}

type KavenegarOptions struct {
	Client        *kavenegar.Kavenegar
	DefaultSender string

	Templates map[string]string
}

func NewKavenegarService(o KavenegarOptions) (SMSService, error) {
	t, err := parseTextTemplates("kavenegar_root", o.Templates)

	return &kavenegarService{
		client:        o.Client,
		defaultSender: o.DefaultSender,
		tpl:           t,
	}, tracer.Trace(err)
}

func (s kavenegarService) Send(_ context.Context, o SMSOptions) error {
	msg, err := s.renderTemplate(o.TemplateName, o.Data)
	if err != nil {
		return tracer.Trace(err)
	}

	if o.Sender == "" {
		o.Sender = s.defaultSender
	}

	_, err = s.client.Message.Send(o.Sender, []string{o.PhoneNumber}, msg, nil)
	return tracer.Trace(err)
}

func (s *kavenegarService) renderTemplate(tplName string, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := s.tpl.ExecuteTemplate(&buf, tplName, data); err != nil {
		return "", tracer.Trace(err)
	}
	return buf.String(), nil
}

// SendVerificationCode ignores the sender param in kavenegar driver.
func (s kavenegarService) SendVerificationCode(_ context.Context, o VerificationOptions) error {
	var lookupParam *kavenegar.VerifyLookupParam
	for _, v := range o.Extra {
		if lp, ok := v.(*kavenegar.VerifyLookupParam); ok {
			lookupParam = lp
		}
	}
	_, err := s.client.Verify.Lookup(o.PhoneNumber, o.TemplateName, o.Code, lookupParam)
	return tracer.Trace(err)
}

var _ SMSService = &kavenegarService{}
