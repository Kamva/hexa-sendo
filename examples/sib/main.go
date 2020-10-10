package main

import (
	"github.com/kamva/gutil"
	sendo "github.com/kamva/hexa-sendo"
	"github.com/kamva/hexa-sendo/sib"
	"os"
)

func main() {
	apiKey := os.Getenv("SIB_API_KEY")

	service := sib.NewEmailService(sib.EmailServiceOptions{
		Client: sib.NewClient(apiKey),
		Templates: map[string]int64{
			"hi": 46,
		},
	})

	err := service.SendSMTP(sendo.SendSMTPEmailOptions{
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
