package main

import (
	"errors"
	"github.com/kamva/sendo"
	"github.com/kavenegar/kavenegar-go"
	"log"
	"os"
	"path"
	"runtime"
)

func main() {
	apiKey := os.Getenv("KAVENEGAR_API_KEY")

	if apiKey == "" {
		log.Fatal(errors.New("provide kavenegar api key please"))
	}
	client := kavenegar.New(apiKey)
	service, err := sendo.NewKavenegarService(sendo.KavenegarOptions{
		Client: client,
		Templates: map[string]string{
			"hi": path.Join(sourcePath(), "templates/hi.tpl"),
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
	return s.SendMessage(sendo.SendSMSOptions{
		TemplateName: "hi",
		PhoneNumber:  "09366579399",
		Data: map[string]interface{}{
			"Name": "مهران",
		},
		Extra: nil,
	})
}

func sendSpeedySMS(s sendo.SMSService) error {
	return s.SendSpeedySMS(sendo.SendSpeedySMSOptions{
		TemplateName: "barekat",
		PhoneNumber:  "09366579399",
		Code:         "K-132443",
		Extra:        nil,
	})
}

func sourcePath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
