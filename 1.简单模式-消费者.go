package main

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	// 1. 连接mq
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

	// 4. 消费消息
	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("消费消息失败：%v", err)
	}

	fmt.Println("等待接收消息...")
	for msg := range msgs {
		fmt.Printf("收到消息：%s\n", msg.Body)
		fmt.Println("等待接收消息...")
		msg.Ack(false)
	}
}
