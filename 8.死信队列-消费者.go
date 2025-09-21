package main

import (
	"go-rabbitmq-demo/core"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	ch := core.InitMQ()
	defer ch.Close()

	exChangeName := "normal_exchange"
	queueName := "normal_queue"
	// 1. 声明交换器
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
		"normal_exchange",    // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("绑定队列失败：%v", err)
	}

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack设为false，手动处理确认
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("消费消息失败：%v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] Received %s", d.Body)
			// 拒绝消息，不重新入队，触发死信机制
			d.Reject(false)
			log.Printf(" [x] Rejected message and sent to DLX")
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
