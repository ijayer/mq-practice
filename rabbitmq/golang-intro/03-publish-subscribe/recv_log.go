/*
 * 说明：
 * 作者：zhe
 * 时间：2019-05-19 11:37 PM
 * 更新：
 */

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ijayer/mq-practice/utils"
	"github.com/streadway/amqp"
)

func main() {
	// rabbitmq server addr
	url := flag.String("url", "amqp://guest:guest@127.0.0.1:5672/", "rabbitmq server address")
	flag.Parse()

	log.Printf("connecting to [%v]\n", *url)
	conn, err := amqp.Dial(*url)
	utils.FatalOnError(err, fmt.Sprintf("failed to connect to rabbitmq[%v]", *url))
	defer conn.Close()
	log.Printf("connected")

	// create channel
	ch, err := conn.Channel()
	utils.FatalOnError(err, "failed to open a channel")
	defer ch.Close()

	// declare exchange
	err = ch.ExchangeDeclare(
		"logs",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	utils.FatalOnError(err, "failed to declare a exchange")

	// declare queue
	queue, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // auto-delete
		true,  // exclusive
		false,
		nil,
	)
	utils.FatalOnError(err, "failed to declare a queue")

	// bind queue with exchange
	err = ch.QueueBind(
		queue.Name,
		"",
		"logs",
		false,
		nil,
	)
	utils.FatalOnError(err, "failed to bind a queue")

	// consume msg
	delivery, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	utils.FatalOnError(err, "failed to register a consume")

	forever := make(chan bool)

	// read and print logs
	go func() {
		for d := range delivery {
			log.Printf("[log] %s\n", string(d.Body))
		}
	}()

	// block
	<-forever
}
