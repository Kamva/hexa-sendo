package sendo

import (
	"bytes"
	"github.com/kamva/tracer"
	"github.com/kavenegar/kavenegar-go"
	"text/template"
)

type KavenegarSender string

// kavenegarService implements the SMSService.
type kavenegarService struct {
	client        *kavenegar.Kavenegar
	defaultSender KavenegarSender

	tpl *template.Template
}

type KavenegarOptions struct {
	Client        *kavenegar.Kavenegar
	DefaultSender KavenegarSender

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

func (s kavenegarService) SendMessage(o SendSMSOptions) error {
	msg, err := s.renderTemplate(o.TemplateName, o.Data)
	if err != nil {
		return tracer.Trace(err)
	}

	sender := s.getSender(o.Extra)

	_, err = s.client.Message.Send(string(sender), []string{o.PhoneNumber}, msg, nil)
	return tracer.Trace(err)
}

func (s *kavenegarService) renderTemplate(tplName string, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := s.tpl.ExecuteTemplate(&buf, tplName, data); err != nil {
		return "", tracer.Trace(err)
	}
	return buf.String(), nil
}

func (s *kavenegarService) getSender(extraOptions []interface{}) KavenegarSender {
	for _, v := range extraOptions {
		if s, ok := v.(KavenegarSender); ok {
			return s
		}
	}
	return s.defaultSender
}

func (s kavenegarService) SendSpeedySMS(o SendSpeedySMSOptions) error {
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
