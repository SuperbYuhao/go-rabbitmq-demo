package main

import (
	"go-rabbitmq-demo/core"
	"log"
	"time"
)

func main() {
	ch := core.InitMQ()
	defer ch.Close()

	queueName := "dl_queue"
	_, err := ch.QueueDeclare(
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

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf(" [x] Received %s", d.Body)
			log.Printf("当前时间: %s", time.Now().Format("2006-01-02 15:04:05"))
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
