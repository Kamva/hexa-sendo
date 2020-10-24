package main

import (
	"github.com/kamva/gutil"
	"github.com/kamva/hexa-sendo"
	"path"
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
	return s.Send(sendo.SMSOptions{
		TemplateName: "hi",
		PhoneNumber:  "09366579399",
		Data: map[string]interface{}{
			"Name": "مهران",
		},
		Extra: nil,
	})
}

func sendSpeedySMS(s sendo.SMSService) error {
	return s.SendVerificationCode(sendo.VerificationOptions{
		TemplateName: "barekat",
		PhoneNumber:  "09366579399",
		Code:         "K-132443",
		Extra:        nil,
	})
}
