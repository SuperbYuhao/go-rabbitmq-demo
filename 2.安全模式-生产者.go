package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	// 1. 加载客户端证书和密钥（双向认证时需要）
	cert, err := tls.LoadX509KeyPair("/Users/yuhao/etc/rabbitmq/ssl/client_certificate.pem", "/Users/yuhao/etc/rabbitmq/ssl/client_key.pem")
	if err != nil {
		log.Fatalf("加载客户端证书失败: %v", err)
	}

	// 2. 加载CA证书（验证服务器证书）
	caCert, err := os.ReadFile("/Users/yuhao/etc/rabbitmq/ssl/ca_certificate.pem")
	if err != nil {
		log.Fatalf("读取CA证书失败: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// 3. 配置TLS
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert}, // 客户端证书（双向认证时需要）
		RootCAs:            caCertPool,              // 信任的CA
		InsecureSkipVerify: false,                   // 必须验证服务器证书
	}

	// 连接 RabbitMQ
	conn, err := amqp091.DialTLS("amqps://admin:password@0.0.0.0:5671/", tlsConfig)
	if err != nil {
		log.Fatalf("无法连接到 RabbitMQ: %v", err)
	}
	defer conn.Close()

	// 2. 创建通道
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("创建通道失败：%v", err)
	}
	defer ch.Close()

	// 3. 声明队列
	q, err := ch.QueueDeclare(
		"easy",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("声明队列失败：%v", err)
	}

	// 4. 发送消息
	body := "Hello Easy Pattern"
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp091.Publishing{
			Body: []byte(body),
		})
	if err != nil {
		log.Fatalf("发送消息失败：%v", err)
	}

	fmt.Printf("发送消息成功：%s \n", body)
}
