package sendo

// Sendo is collection of send services.
type Sendo interface {
	Email() EmailService
	SMS() SMSService
}

type sendoList struct {
	email EmailService
	sms   SMSService
}

type Builder struct {
	services sendoList
}

func (b *Builder) WithEmail(e EmailService) {
	b.services.email = e
}

func (b *Builder) WithSMS(sms SMSService) {
	b.services.sms = sms
}

func (b *Builder) Build() Sendo {
	return &sendo{
		services: b.services,
	}
}

type sendo struct {
	services sendoList
}

func (s *sendo) Email() EmailService {
	return s.services.email
}

func (s *sendo) SMS() SMSService {
	return s.services.sms
}

var _ Sendo = &sendo{}
