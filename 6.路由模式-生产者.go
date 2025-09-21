package main

import (
	"fmt"
	"go-rabbitmq-demo/core"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	ch := core.InitMQ()
	defer ch.Close()

	exChangeName := "router"
	// 1. 声明 fanout 类型交换器
	err := ch.ExchangeDeclare(
		exChangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("声明交换器失败：%v", err)
	}

	// 2. 将消息发给交换器
	for i := 0; i < 20; i++ {
		// 发送消息
		body := fmt.Sprintf("msg-%d", i)
		err = ch.Publish(
			exChangeName, //  交换器名称
			"node001",
			false,
			false,
			amqp091.Publishing{
				Body: []byte(body),
			})
		if err != nil {
			fmt.Printf("发送失败: %v", err)
			return
		}
		fmt.Printf("发送消息: %s\n", body)
	}
	for i := 0; i < 10; i++ {
		// 发送消息
		body := fmt.Sprintf("msg-%d", i)
		err = ch.Publish(
			exChangeName, //  交换器名称
			"node002",
			false,
			false,
			amqp091.Publishing{
				Body: []byte(body),
			})
		if err != nil {
			fmt.Printf("发送失败: %v", err)
			return
		}
		fmt.Printf("发送消息: %s\n", body)
	}
}
