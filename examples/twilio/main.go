package main

import (
	"errors"
	"github.com/kamva/gutil"
	"github.com/kamva/hexa-sendo"
	"log"
	"os"
	"path"
)

func main() {
	sid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")

	if sid == "" || authToken == "" {
		log.Fatal(errors.New("provide twilio sid & auth_token"))
	}

	service, err := sendo.NewTwilioService(sendo.TwilioOptions{
		AccountSID: sid,
		AuthToken:  authToken,
		Templates: map[string]string{
			"verify": path.Join(gutil.SourcePath(), "templates/verify.tpl"),
		},
	})
	if err != nil {
		panic(err)
	}
	if err := sendVerificationSMS(service); err != nil {
		panic(err)
	}
}

func sendVerificationSMS(s sendo.SMSService) error {
	return s.SendVerificationCode(sendo.VerificationOptions{
		TemplateName: "verify",
		PhoneNumber:  "+989130022039",
		Code:         "132443",
		Sender:       os.Getenv("TWILIO_SENDER"),
		Extra:        nil,
	})
}
