package main

import (
	"fmt"
	"go-rabbitmq-demo/core"
	"log"
	"os"
)

func main() {
	ch := core.InitMQ()
	defer ch.Close()

	exChangeName := "topic"
	// 1. 声明交换器
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

	// 2. 创建临时队列
	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("创建临时队列失败：%v", err)
	}

	// 获取命令行参数作为绑定键
	if len(os.Args) < 2 {
		log.Printf("Usage: %s [binding_key]...", os.Args[0])
		os.Exit(0)
	}
	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s",
			q.Name, "topic", s)
		err = ch.QueueBind(
			q.Name,       // queue name
			s,            // routing key
			exChangeName, // exchange
			false,
			nil)
		if err != nil {
			log.Fatalf("绑定队列失败：%v", err)
		}
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

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s: %s", d.RoutingKey, d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
