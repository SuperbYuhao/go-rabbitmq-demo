package main

import (
	"fmt"
	"go-rabbitmq-demo/core"
	"log"
	"os"
	"strings"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	ch := core.InitMQ()
	defer ch.Close()

	exChangeName := "topic"
	// 1. 声明 fanout 类型交换器
	err := ch.ExchangeDeclare(
		exChangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("声明交换器失败：%v", err)
	}

	// 设置路由键，默认为"anonymous.info"
	body := bodyFrom(os.Args)
	routingKey := severityFrom(os.Args)

	// 将消息发给交换器
	err = ch.Publish(
		exChangeName, //  交换器名称
		routingKey,   // 路由键
		false,        // mandatory
		false,        // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		fmt.Printf("发送失败: %v", err)
		return
	}
	fmt.Printf("发送消息至 %s : %s\n", routingKey, body)
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 3) || os.Args[2] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[2:], " ")
	}
	return s
}

func severityFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "anonymous.info"
	} else {
		s = os.Args[1]
	}
	return s
}
