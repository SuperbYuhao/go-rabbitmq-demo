package main

import (
	"fmt"
	"go-rabbitmq-demo/core"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	ch := core.InitMQ()

	q, err := ch.QueueDeclare("work-queue", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("声明队列失败：%v", err)
	}

	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("工作队列：%d", i)
		err = ch.Publish("", q.Name, false, false, amqp091.Publishing{
			Body: []byte(msg)},
		)
		if err != nil {
			fmt.Printf("生产消息失败：%v \n", err)
			return
		}
		fmt.Printf("消息发送成功: %s\n", msg)
	}

}
