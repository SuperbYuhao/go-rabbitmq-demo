package main

import (
	"fmt"
	"go-rabbitmq-demo/core"
	"log"
)

func main() {
	ch := core.InitMQ()
	defer ch.Close()

	q, err := ch.QueueDeclare("work-queue", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("声明队列失败：%v", err)
	}

	// 消费者1
	//err = ch.Qos(
	//	1,     // prefetchCount: 未确认消息上限
	//	0,     // prefetchSize: 不限制消息大小
	//	false, // global: 仅对当前消费者生效
	//)
	// 消费者2
	err = ch.Qos(
		2,     // prefetchCount: 未确认消息上限
		0,     // prefetchSize: 不限制消息大小
		false, // global: 仅对当前消费者生效
	)

	// 4. 消费消息
	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
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
		fmt.Println("等待接收消息...")
		msg.Ack(false)
	}
}
