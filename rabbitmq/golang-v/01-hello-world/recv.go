/*
 * 说明：
 * 作者：zhe
 * 时间：2019-03-28 9:22 PM
 * 更新：
 */

package main

import (
	"flag"
	"log"

	"github.com/streadway/amqp"
	"github.com/zhezh09/mq-practice/utils"
)

func main() {
	// rabbitmq server addr
	url := flag.String("url", "amqp://guest:guest@10.0.0.69:5672/", "rabbitmq server address")
	flag.Parse()

	// 连接到 RabbitMQ 服务器，即：建立 Socket 连接，处理
	// 协议转换、版本对接以及一些登陆授权问题 For Us.
	log.Printf("connecting to [%s]\n", *url)
	conn, err := amqp.Dial(*url)
	utils.FatalOnError(err, "failed to connect to RabbitMQ")
	log.Println("consumer connected")

	// 确保连接被关闭
	defer func() { _ = conn.Close() }()

	// 创建一个通道，然后通过这个 ch 来实现我们的相关 API
	ch, err := conn.Channel()
	utils.FatalOnError(err, "failed to open a channel")
	defer func() { _ = ch.Close() }()

	// 创建一个队列，用来存储、转发消息
	// 消费者只需要从这个 queue 读取消息，就完成了 Delivery
	//
	// Note：Queue 的 name 必须和生产者定义保持一致。这样才
	// 能实现发送/接收相匹配
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable 持久化
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	utils.FatalOnError(err, "failed to declare a queue")

	// 注册一个消费者, 队列中的消息将被传送到 ‘chan Delivery’ 这个通道中
	consumeCh, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	utils.FatalOnError(err, "failed to register a consumer")

	forever := make(chan struct{})

	// 用一个 Go 程，持续从 consumeCh 中读取 queue 传送的消息
	go func() {
		var d amqp.Delivery
		for d = range consumeCh {
			log.Printf("received a message: %s\n", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	// 2019/03/29 14:13:14 connecting to [amqp://guest:guest@10.0.0.69:5672/]
	// 2019/03/29 14:13:14 consumer connected
	// 2019/03/29 14:13:14  [*] Waiting for messages. To exit press CTRL+C
	// 2019/03/29 14:13:14 received a message: Hello
}
