package main

import (
	"fmt"
	"user-microservice/kafka"
)

func (app *application) signupEmailKafkaProducer(newUserEmail string) {
	app.logger.Printf("Producing event to send signup email to user: %s", newUserEmail)

	go kafka.SignupProducer(app.producers.signup, newUserEmail)
	app.logger.Println("Fired a goroutine to kafka signup producer")
}

func (app *application) resetPasswordKafkaProducer(email string, verificationCode int, ipAddress string) {
	fmt.Printf("Producing reset password event for email: %s, verification code is %d, ip addr: %s",
		email,
		verificationCode,
		ipAddress)

	go kafka.ResetPasswordProducer(app.producers.reset, email, verificationCode, ipAddress)
	app.logger.Println("Fired a goroutine to kafka reset password producer")
}
