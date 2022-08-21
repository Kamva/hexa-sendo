package main

import (
	"context"
	"errors"
	"log"
	"os"
	"path"

	"github.com/kamva/gutil"
	sendo "github.com/kamva/hexa-sendo"
	"github.com/kamva/hexa/hexahttp"
)

func main() {
	apiKey := os.Getenv("MEDIANA_API_KEY")

	if apiKey == "" {
		log.Fatal(errors.New("provide mediana api key please"))
	}

	cli, err := hexahttp.NewClient("http://apixxx.xxx.xx/api/v1", hexahttp.LogModeNone)

	service, err := sendo.NewMedianaService(sendo.MedianaOptions{
		ApiClient: cli,
		Templates: map[string]string{
			"hi": path.Join(gutil.SourcePath(), "templates/hi.tpl"),
		},
		Token:         apiKey,
		DefaultSender: "+989xxxxxx",
	})
	if err != nil {
		panic(err)
	}
	if err := sendSMS(service); err != nil {
		panic(err)
	}
}

func sendSMS(s sendo.SMSService) error {
	return s.Send(context.Background(), sendo.SMSOptions{
		TemplateName: "hi",
		PhoneNumber:  "98912xxxxxxxx",
		Data: map[string]interface{}{
			"Name": "مهران",
		},
	})
}
