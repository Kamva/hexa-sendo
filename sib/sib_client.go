package sib

import (
	"encoding/json"
	"fmt"
	"github.com/kamva/gutil"
	"github.com/kamva/hexa/hexahttp"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
)

type SendinblueClient struct {
	auth   hexahttp.RequestOption
	client *hexahttp.Client
}

type SendSMTPEmailResponse struct {
	MessageID string `json:"messageId"`
}

func (c *SendinblueClient) SendSMTPEmail(p SendSMTPEmailParams) (*SendSMTPEmailResponse, error) {
	params := p.RequestParams()
	resp, err := c.client.PostJsonFormWithOptions("smtp/email", params, c.auth)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	if resp.StatusCode != 201 {
		err := fmt.Errorf("error with status code %v", resp.StatusCode)

		hlog.Error("error on send SMTP email",
			hlog.Err(err),
			hlog.String("params", string(gutil.Must(json.Marshal(params)).([]byte))),
			hlog.String("resp", fmt.Sprintf("%+v", resp)),
			hlog.Int("status_code", resp.StatusCode),
		)
		return nil, tracer.Trace(err)
	}
	b, err := hexahttp.Bytes(resp, err)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	var data SendSMTPEmailResponse
	return &data, tracer.Trace(json.Unmarshal(b, &resp))
}

func NewClient(apiKey string) *SendinblueClient {
	return &SendinblueClient{
		client: hexahttp.NewClient(gutil.NewString("https://api.sendinblue.com/v3")),
		auth:   hexahttp.AuthenticateHeader("api-key", "", apiKey),
	}
}
