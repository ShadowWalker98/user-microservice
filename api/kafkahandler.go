package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	kafka2 "user-microservice/kafka"
)

func (app *application) signupEmailKafkaProducer(newUserEmail string) {
	app.logger.Printf("Producing event to send signup email to user: %s", newUserEmail)
	signupProducer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":       "192.168.0.9:9092",
		"socket.keepalive.enable": true,
		"log.connection.close":    false,
	})
	if err != nil {
		app.logger.Println("error while initialising kafka connection for signup producer")
	}
	go kafka2.SignupProducer(signupProducer, newUserEmail)
	app.logger.Println("Fired a goroutine to kafka signup producer")
}

func (app *application) resetPasswordKafkaProducer(email string, verificationCode int, ipAddress string) {
	fmt.Printf("Producing reset password event for email: %s, verification code is %d, ip addr: %s",
		email,
		verificationCode,
		ipAddress)
	resetProducer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":       "192.168.0.9:9092",
		"socket.keepalive.enable": true,
		"log.connection.close":    false,
	})

	if err != nil {
		app.logger.Println("error while initialising kafka connection for reset password producer")
	}
	go kafka2.ResetPasswordProducer(resetProducer, email, verificationCode, ipAddress)
	app.logger.Println("Fired a goroutine to kafka reset password producer")
}
