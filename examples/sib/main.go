package main

import (
	"context"
	"os"

	"github.com/kamva/gutil"
	sendo "github.com/kamva/hexa-sendo"
	"github.com/kamva/hexa-sendo/sib"
	"github.com/kamva/hexa/hexahttp"
)

func main() {
	apiKey := os.Getenv("SIB_API_KEY")
	cli, err := sib.NewClientWithDefaults(apiKey, hexahttp.LogModeNone)
	gutil.PanicErr(err)

	service := sib.NewEmailService(sib.EmailServiceOptions{
		Client: cli,
		Templates: map[string]int64{
			"hi": 46,
		},
	})

	err = service.SendSMTP(context.Background(), sendo.SendSMTPEmailOptions{
		To: []sendo.EmailTo{
			{
				Name:  "mehran prs",
				Email: "mehr.prs@gmail.com",
			},
		},
		Subject:      gutil.NewString("Hi from sendo"),
		TemplateName: "hi",
	})
	gutil.PanicErr(err)
}
