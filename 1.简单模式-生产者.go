package main

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	// 1. 连接RabbitMQ
	conn, err := amqp091.Dial("amqp://admin:password@localhost:5672/")
	if err != nil {
		log.Fatalf("连接mq失败：%v", err)
	}
	defer conn.Close()

	// 2. 创建通道
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("创建通道失败：%v", err)
	}
	defer ch.Close()

	// 3. 声明队列
	q, err := ch.QueueDeclare(
		"easy",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("声明队列失败：%v", err)
	}

	// 4. 发送消息
	body := "Hello Easy Pattern"
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp091.Publishing{
			Body: []byte(body),
		})
	if err != nil {
		log.Fatalf("发送消息失败：%v", err)
	}

	fmt.Printf("发送消息成功：%s \n", body)
}
