package sib

import (
	"errors"
	"github.com/kamva/gutil"
	"github.com/kamva/hexa-sendo"
	"github.com/kamva/tracer"
)

type sendinBlueEmailService struct {
	client    *SendinblueClient
	templates map[string]int64
}

func (s *sendinBlueEmailService) SendSMTP(o sendo.SendSMTPEmailOptions) error {
	tplID, ok := s.templates[o.TemplateName]
	if !ok {
		return tracer.Trace(errors.New("can not send email using sib, template id not found"))
	}
	_, err := s.client.SendSMTPEmail(SendSMTPEmailParams{
		Sender:     o.Sender,
		From:       o.From,
		To:         o.To,
		CC:         o.CC,
		ReplyTo:    o.ReplyTo,
		Subject:    o.Subject,
		TemplateID: tplID,
		Params:     gutil.StructToMap(o.Data),
	})

	return tracer.Trace(err)
}

type EmailServiceOptions struct {
	Client    *SendinblueClient
	Templates map[string]int64 // mapping: {template name} => {sib template id}
}

func NewEmailService(o EmailServiceOptions) sendo.EmailService {
	return &sendinBlueEmailService{
		client:    o.Client,
		templates: o.Templates,
	}
}

var _ sendo.EmailService = &sendinBlueEmailService{}
