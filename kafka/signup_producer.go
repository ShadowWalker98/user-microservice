package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type UserInfo struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type ResetPasswordInfo struct {
	Email            string `json:"email"`
	VerificationCode int    `json:"verification_code"`
	IPAddress        string `json:"ip_address"`
}

const (
	KafkaSignupTopic = "workouts-v1-users-test"
	KafkaResetTopic  = "workouts-v1-users-password-reset-test"
)

func SignupProducer(p *kafka.Producer, email string, userId int) {

	fmt.Println("Creating a user")

	topic := KafkaSignupTopic
	user := UserInfo{
		ID:    userId,
		Email: email,
	}

	value, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	fmt.Println("Pushing user onto kafka queue")

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: value,
	}, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Waiting for kafka to push message to the topic!")
	p.Flush(5000)
}

func ResetPasswordProducer(p *kafka.Producer, email string, verificationCode int, ipAddress string) {
	fmt.Printf("Resetting password for user %s, verification code generated: %d and originating IP: %s",
		email,
		verificationCode,
		ipAddress)

	topic := KafkaResetTopic

	resetPasswordInfo := ResetPasswordInfo{
		Email:            email,
		VerificationCode: verificationCode,
		IPAddress:        ipAddress,
	}

	value, err := json.Marshal(resetPasswordInfo)
	if err != nil {
		panic(err)
	}

	fmt.Println("Pushing reset password info onto kafka queue")

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: value,
	}, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Waiting for kafka to push message to the topic!")
	p.Flush(5000)
}
