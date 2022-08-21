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
	resp, err := c.client.PostJsonForm("smtp/email", params, c.auth)
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

// NewClient creates a new Sendinblue client.
// The hexahttp Client param must have the base sendinblue API url as a base url.
// e.g., hexahttp.NewClient("https://api.sendinblue.com/v3")
func NewClient(cli *hexahttp.Client, apiKey string) *SendinblueClient {
	return &SendinblueClient{
		client: cli,
		auth:   hexahttp.AuthenticateHeader("api-key", "", apiKey),
	}
}

func NewClientWithDefaults(apiKey string, httpClientLogMode uint) (*SendinblueClient, error) {
	cli, err := hexahttp.NewClient("https://api.sendinblue.com/v3", httpClientLogMode)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	return NewClient(cli, apiKey), nil
}
