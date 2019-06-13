/*
 * 说明：
 * 作者：zhe
 * 时间：2019-05-19 11:36 PM
 * 更新：
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
	log.Println("connected")

	// 确保连接被关闭
	defer conn.Close()

	// 创建一个通道，然后通过这个 ch 来实现我们的相关 API
	ch, err := conn.Channel()
	utils.FatalOnError(err, "failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto delete
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	utils.FatalOnError(err, "failed to declare an exchange")

	body := utils.BodyFrom(os.Args)
	err = ch.Publish(
		"logs",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	utils.FatalOnError(err, "failed to publish a msg")

	log.Printf("[X] sent: %s \n", body)
}
