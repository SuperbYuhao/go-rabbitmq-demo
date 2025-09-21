package main

import (
	"go-rabbitmq-demo/core"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	ch := core.InitMQ()
	defer ch.Close()

	exChangeName := "dlx_exchange"
	queueName := "dl_queue"
	delayQueueName := "delay_queue"

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

	err = ch.QueueBind(
		queueName,    // queue name
		"dl_key",     // routing key
		exChangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("绑定队列失败：%v", err)
	}

	_, err = ch.QueueDeclare(
		delayQueueName, // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		amqp091.Table{
			"x-dead-letter-exchange":    exChangeName,
			"x-dead-letter-routing-key": "dl_key",
			"x-message-ttl":             5000, // 5秒后消息过期
		},
	)
	if err != nil {
		log.Fatalf("创建延迟队列失败：%v", err)
	}

	// 发送消息到延迟队列
	body := "Hello, delayed message!"
	err = ch.Publish(
		"",             // exchange
		delayQueueName, // routing key
		false,          // mandatory
		false,          // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Fatalf("生产消息失败：%v", err)
	}

	log.Printf(" [x] Sent %s", body)
	log.Println("消息已发送到延迟队列，5秒后将进入死信队列")

	time.Sleep(2 * time.Second) // 等待消息发送完成
}
