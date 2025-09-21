package main

import (
	"go-rabbitmq-demo/core"
	"log"
)

func main() {
	ch := core.InitMQ()
	defer ch.Close()

	exChangeName := "dlx_exchange"
	queueName := "dlx_queue"
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

	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("创建队列失败：%v", err)
	}

	// 绑定死信队列到死信交换器
	err = ch.QueueBind(
		queueName,    // queue name
		"dlx_key",    // routing key
		exChangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("绑定队列失败：%v", err)
	}

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack设为false，手动处理确认
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
