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

	exChangeName := "pub-sub-cache"
	// 1. 声明 fanout 类型交换器
	err := ch.ExchangeDeclare(
		exChangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("声明交换器失败：%v", err)
	}

	// 声明队列1
	q1, err := ch.QueueDeclare(
		"pub_sub_queue_2",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("声明队列1失败：%v", err)
	}
	// 声明队列2
	q2, err := ch.QueueDeclare(
		"pub_sub_queue_2",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("声明队列2失败：%v", err)
	}

	// 绑定队列1
	err = ch.QueueBind(
		q1.Name,      // 队列名称
		"",           // 路由键（fanout 类型忽略）
		exChangeName, // 交换器名称
		false,        // 不等待服务器确认
		nil,          // 额外参数
	)
	// 绑定队列2
	err = ch.QueueBind(
		q2.Name,      // 队列名称
		"",           // 路由键（fanout 类型忽略）
		exChangeName, // 交换器名称
		false,        // 不等待服务器确认
		nil,          // 额外参数
	)
	if err != nil {
		log.Fatalf("队列绑定交换器失败 %s", err)
	}

	// 最后. 将消息发给交换器
	for i := 0; i < 20; i++ {
		// 发送消息
		body := fmt.Sprintf("msg-%d", i)
		err = ch.Publish(
			exChangeName, //  交换器名称
			"",           // 路由键（fanout 类型忽略）
			false,        // mandatory
			false,        // immediate
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
