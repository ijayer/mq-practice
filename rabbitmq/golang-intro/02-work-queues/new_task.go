/*
 * 说明：RabbitMQ: Work Queues
 * 作者：zhe
 * 时间：2019-03-28 9:22 PM
 * 更新：RabbitMQ API: https://godoc.org/github.com/streadway/amqp
 */

package main

import (
	"flag"
	"log"
	"os"

	"github.com/ijayer/mq-practice/utils"
	"github.com/streadway/amqp"
)

func main() {
	// rabbitmq server addr
	url := flag.String("url", "amqp://guest:guest@127.0.0.1:5672/", "rabbitmq server address")
	flag.Parse()

	// 连接到 RabbitMQ 服务器，即：建立 Socket 连接，处理
	// 协议转换、版本对接以及一些登陆授权问题 For Us.
	log.Printf("connecting to [%s]\n", *url)
	conn, err := amqp.Dial(*url)
	utils.FatalOnError(err, "failed to connect to RabbitMQ")
	log.Println("producer connected")

	// 确保连接被关闭
	defer func() { _ = conn.Close() }()

	// 创建一个通道，然后通过这个 ch 来实现我们的相关 API
	ch, err := conn.Channel()
	utils.FatalOnError(err, "failed to open a channel")
	defer func() { _ = ch.Close() }()

	// 创建一个队列，用来存储、转发消息
	// 生产者只需要将消息写入这个 queue 就完成了 Publishing
	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable 持久化
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	utils.FatalOnError(err, "failed to declare a queue")

	// 发送消息，消息内容从命令行获取
	body := utils.BodyFrom(os.Args)

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	utils.FatalOnError(err, "failed to publish a message")
	log.Printf(" [x] Sent %s", body)

	// 2019/03/29 10:43:48 connecting to [amqp://guest:guest@127.0.0.1:5672/]
	// 2019/03/29 10:43:48 connected
	// 2019/03/29 10:43:48  [x] Sent Hello
}
