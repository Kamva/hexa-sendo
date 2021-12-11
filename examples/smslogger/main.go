package main

import (
	"context"
	"path"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa-sendo"
)

func main() {
	service, err := sendo.NewSMSLoggerService(sendo.SMSLoggerOptions{
		Templates: map[string]string{
			"hi": path.Join(gutil.SourcePath(), "templates/hi.tpl"),
		},
	})
	if err != nil {
		panic(err)
	}
	if err := sendSMS(service); err != nil {
		panic(err)
	}
	if err := sendSpeedySMS(service); err != nil {
		panic(err)
	}
}

func sendSMS(s sendo.SMSService) error {
	return s.Send(context.Background(), sendo.SMSOptions{
		TemplateName: "hi",
		PhoneNumber:  "09366579399",
		Data: map[string]interface{}{
			"Name": "مهران",
		},
		Extra: nil,
	})
}

func sendSpeedySMS(s sendo.SMSService) error {
	return s.SendVerificationCode(context.Background(), sendo.VerificationOptions{
		TemplateName: "barekat",
		PhoneNumber:  "09366579399",
		Code:         "K-132443",
		Extra:        nil,
	})
}
