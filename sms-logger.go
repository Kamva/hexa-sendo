package sendo

import (
	"bytes"
	"fmt"
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"text/template"
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

func (s smsLoggerService) Send(o SMSOptions) error {
	msg, err := s.renderTemplate(o.TemplateName, o.Data)
	if err != nil {
		return tracer.Trace(err)
	}
	hlog.WithFields(gutil.MapToKeyValue(hexa.Map{
		"to":            o.PhoneNumber,
		"message":       msg,
		"template_name": o.TemplateName,
		"data":          fmt.Sprintf("%+v", o.Data),
	})).Message("send sms")
	return nil
}

func (s *smsLoggerService) renderTemplate(tplName string, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := s.tpl.ExecuteTemplate(&buf, tplName, data); err != nil {
		return "", tracer.Trace(err)
	}
	return buf.String(), nil
}

func (s smsLoggerService) SendVerificationCode(o VerificationOptions) error {
	hlog.WithFields("to", o.PhoneNumber, "code", o.Code, "template_name", o.TemplateName).Message("send sms")
	return nil
}

var _ SMSService = &smsLoggerService{}
