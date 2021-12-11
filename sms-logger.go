package sendo

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
)

type SMSLoggerOptions struct {
	Templates map[string]string
}

// smsLogger implements the SMSService.
type smsLoggerService struct {
	tpl *template.Template
}

// NewSMSLoggerService returns new instance of the sms logger.
func NewSMSLoggerService(o SMSLoggerOptions) (SMSService, error) {
	t, err := parseTextTemplates("smsLogger_root", o.Templates)

	return &smsLoggerService{
		tpl: t,
	}, tracer.Trace(err)
}

func (s smsLoggerService) Send(_ context.Context, o SMSOptions) error {
	msg, err := s.renderTemplate(o.TemplateName, o.Data)
	if err != nil {
		return tracer.Trace(err)
	}
	hlog.Message("send sms",
		hlog.String("to", o.PhoneNumber),
		hlog.String("message", msg),
		hlog.String("template_name", o.TemplateName),
		hlog.String("data", fmt.Sprintf("%+v", o.Data)),
	)
	return nil
}

func (s *smsLoggerService) renderTemplate(tplName string, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := s.tpl.ExecuteTemplate(&buf, tplName, data); err != nil {
		return "", tracer.Trace(err)
	}
	return buf.String(), nil
}

func (s smsLoggerService) SendVerificationCode(_ context.Context, o VerificationOptions) error {
	hlog.Message("send sms",
		hlog.String("to", o.PhoneNumber),
		hlog.String("code", o.Code),
		hlog.String("template_name", o.TemplateName))
	return nil
}

var _ SMSService = &smsLoggerService{}
