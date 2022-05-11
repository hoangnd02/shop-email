package services

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/smtp"

	"github.com/hoanggggg5/shopemail/config"
	"github.com/twmb/franz-go/pkg/kgo"
)

type EngineMailerPayload struct {
	Key    string                 `json:"key"`
	To     string                 `json:"to"`
	Record map[string]interface{} `json:"record"`
}

type SendEmail struct{}

func NewSendEmail() *SendEmail {
	return &SendEmail{}
}

func (SendEmail) Process() {
	cl, err := kgo.NewClient(
		kgo.SeedBrokers("localhost:9092"),
		kgo.ConsumerGroup("test"),
		kgo.ConsumeTopics("mailer"),
		kgo.AllowAutoTopicCreation(),
		kgo.DisableAutoCommit(),
	)
	if err != nil {
		panic(err)
	}

	for {
		fetches := cl.PollFetches(context.Background())
		// We can iterate through a record iterator...
		records := fetches.Records()
		for _, r := range records {
			cl.CommitRecords(context.Background(), r)

			mailerPayload := EngineMailerPayload{}
			err = json.Unmarshal(r.Value, &mailerPayload)
			if err != nil {
				panic("could not parse mailerPayload" + err.Error())
			}

			for _, event := range config.Mailer.Events {
				if event.Key == mailerPayload.Key {
					template_path := event.Template
					content_tpl, err := template.ParseFiles(template_path)
					if err != nil {
						log.Println(err)
					}

					content_buf := bytes.Buffer{}
					content_tpl.Execute(&content_buf, mailerPayload)

					msg := []byte("\r\n" + content_buf.String())

					log.Println("msg: ", string(msg), mailerPayload.To)
					SendEmailService(msg, mailerPayload.To)
					break
				}
			}
		}
	}
}

func SendEmailService(message []byte, toAddress string) (response bool, err error) {
	fromAddress := "nghiemdanghoang@gmail.com"
	fromEmailPassword := "14124869"
	smtpServer := "smtp.gmail.com"
	smptPort := "587"

	var auth = smtp.PlainAuth("", fromAddress, fromEmailPassword, smtpServer)
	if err := smtp.SendMail(smtpServer+":"+smptPort, auth, fromAddress, []string{toAddress}, message); err == nil {
		log.Println("Send email successfully")
		return true, nil
	} else {
		log.Println("Send email error")
		return false, err
	}
}
