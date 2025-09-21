package core

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func InitMQ() (ch *amqp091.Channel) {
	// 1. 连接mq
	conn, err := amqp091.Dial("amqp://admin:password@localhost:5672/")
	if err != nil {
		log.Fatalf("连接mq失败：%v", err)
	}

	// 2. 创建通道
	ch, err = conn.Channel()
	if err != nil {
		log.Fatalf("创建通道失败：%v", err)
	}

	return ch
}
