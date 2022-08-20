package sendo

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hexahttp"
)

// medianaService implements the SMSService.
type medianaService struct {
	client        *hexahttp.Client
	defaultSender string
	token         string
}

type MedianaOptions struct {
	APIUrl        string
	Token         string
	DefaultSender string
}

func NewMedianaService(o MedianaOptions) (SMSService, error) {

	return &medianaService{
		client:        hexahttp.NewClient(&o.APIUrl),
		defaultSender: o.DefaultSender,
		token:         o.Token,
	}, nil
}

type medianaMessage struct {
	Message string
}

func (s medianaService) Send(_ context.Context, o SMSOptions) error {
	var authorizationHeader = hexahttp.AuthenticateHeader("apikey", "", s.token)
	var recipients []string
	recipients = append(recipients, "+"+o.PhoneNumber)
	sender := o.Sender
	if sender == "" {
		sender = s.defaultSender
	}

	dataJsonStr, err := json.Marshal(o.Data)
	if err != nil {
		return err
	}
	var medianaMsg medianaMessage
	if err = json.Unmarshal(dataJsonStr, &medianaMsg); err != nil {
		return err
	}
	resp, err := s.client.PostJsonFormWithOptions("sms/send/webservice/single", hexa.Map{
		"recipient": recipients,
		"sender":    sender,
		"message":   medianaMsg.Message,
	}, authorizationHeader)
	if err != nil {
		return err
	}
	return hexahttp.ResponseError(resp)
}

func (s *medianaService) renderTemplate(tplName string, data interface{}) (string, error) {
	var buf bytes.Buffer
	return buf.String(), nil
}

// SendVerificationCode ignores the sender param in mediana driver.
func (s medianaService) SendVerificationCode(_a context.Context, o VerificationOptions) error {

	return nil
}

var _ SMSService = &kavenegarService{}
