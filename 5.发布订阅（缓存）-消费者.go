package main

import (
	"fmt"
	"go-rabbitmq-demo/core"
	"log"
)

func main() {
	ch := core.InitMQ()
	defer ch.Close()

	exChangeName := "pub-sub-cache"
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

	// 2. 消费消息
	msgs1, err := ch.Consume(
		"pub_sub_queue_1",
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
	for msg := range msgs1 {
		fmt.Printf("收到消息：%s\n", msg.Body)
	}

	// 2. 消费消息
	msgs, err := ch.Consume(
		"pub_sub_queue_2",
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
