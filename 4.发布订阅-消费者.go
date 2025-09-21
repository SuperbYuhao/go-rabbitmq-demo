package main

import (
	"fmt"
	"go-rabbitmq-demo/core"
	"log"
)

func main() {
	ch := core.InitMQ()
	defer ch.Close()

	exChangeName := "pub-sub"
	// 1. 声明交换器
	err := ch.ExchangeDeclare(
		exChangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("声明交换器失败：%v", err)
	}

	// 2. 创建临时队列
	q, err := ch.QueueDeclare(
		"",
		false,
		true,
		true,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("创建临时队列失败：%v", err)
	}

	// 3. 将队列绑定到交换器
	err = ch.QueueBind(
		q.Name,
		"",
		exChangeName,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("队列绑定到交换器失败：%v", err)
	}

	// 4. 消费消息
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
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
	}
}
