package main

import (
	"go-rabbitmq-demo/core"
	"log"
	"os"
	"strings"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	ch := core.InitMQ()
	defer ch.Close()

	exChangeName := "normal_exchange"
	queueName := "normal_queue"
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

	args := amqp091.Table{
		"x-dead-letter-exchange":    "dlx_exchange",
		"x-dead-letter-routing-key": "dlx_key",
	}
	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		args,      // arguments
	)
	if err != nil {
		log.Fatalf("创建队列失败：%v", err)
	}

	err = ch.QueueBind(
		queueName,            // queue name
		"normal_routing_key", // routing key
		exChangeName,         // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("绑定队列失败：%v", err)
	}

	body := bodyFrom8(os.Args)
	err = ch.Publish(
		exChangeName,         // exchange
		"normal_routing_key", // routing key
		false,                // mandatory
		false,                // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Fatalf("发送消息失败：%v", err)
	}

	log.Printf(" [x] Sent %s", body)
}

func bodyFrom8(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
