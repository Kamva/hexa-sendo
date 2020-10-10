package sib

import (
	"github.com/kamva/hexa"
	sendo "github.com/kamva/hexa-sendo"
)

type SendSMTPEmailParams struct {
	Sender     *sendo.EmailSender     // optional, driver should support default sender per for each template
	From       *sendo.EmailFrom       // From for some drivers can be optional, for some other required.
	To         []sendo.EmailTo        // required
	CC         []sendo.EmailCC        // optional
	ReplyTo    *sendo.EmailReplyTo    // optional
	Subject    *string                // optional, drivers must support default subject for each template.
	TemplateID int64                  // required
	Params     map[string]interface{} // optional
	Extra      []interface{}          // optional extra options.
}

func (p SendSMTPEmailParams) ToParam() []hexa.Map {
	var to = make([]hexa.Map, len(p.To))
	for k, v := range p.To {
		to[k] = hexa.Map{"name": v.Name, "email": v.Email}
	}
	return to
}

func (p SendSMTPEmailParams) CCParam() []hexa.Map {
	if len(p.CC) == 0 {
		return nil
	}
	cc := make([]hexa.Map, len(p.CC))
	for k, v := range p.CC {
		cc[k] = hexa.Map{"name": v.Name, "email": v.Email}
	}
	return cc
}
func (p SendSMTPEmailParams) SenderParam() hexa.Map {
	if p.Sender == nil {
		return nil
	}
	sender := hexa.Map{
		"email": p.Sender.Email,
	}
	if p.Sender.Name != "" {
		sender["name"] = p.Sender.Name
	}
	return sender
}

func (p SendSMTPEmailParams) ReplyToParam() hexa.Map {
	if p.ReplyTo == nil {
		return nil
	}
	return hexa.Map{
		"name":  p.ReplyTo.Name,
		"email": p.ReplyTo.Email,
	}
}

func (p SendSMTPEmailParams) RequestParams() hexa.Map {
	data := hexa.Map{
		"templateId": p.TemplateID,
		"to":         p.ToParam(),
	}

	if len(p.Params) != 0 {
		data["params"] = p.Params
	}

	if cc := p.CCParam(); cc != nil {
		data["cc"] = cc
	}
	if sender := p.SenderParam(); sender != nil {
		data["sender"] = sender
	}

	if p.Subject != nil {
		data["subject"] = p.Subject
	}

	if replyTo := p.ReplyToParam(); replyTo != nil {
		data["replyTo"] = replyTo
	}

	return data
}
