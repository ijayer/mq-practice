/*
 * 说明：
 * 作者：jayer
 * 时间：2019-06-02 10:54 PM
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

	log.Printf("connecting to rabbitmq[%v]\n", *url)
	// dial with rabbitmq
	conn, err := amqp.Dial(*url)
	utils.FatalOnError(err, "failed to connect to rabbitmq")
	log.Printf("connected")

	// create channel
	ch, err := conn.Channel()
	utils.FatalOnError(err, "failed to create a channel")

	// declare exchange
	xType := "direct"
	xName := "logs_direct"
	err = ch.ExchangeDeclare(xName, xType, true, false, false, false, nil)
	utils.FatalOnError(err, "failed to declare an exchange")

	// declare queue
	queue, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	utils.FatalOnError(err, "failed to create a queue")

	if len(os.Args) < 2 {
		log.Printf("Usage: %v [info] [warning] [error]", os.Args[0])
		os.Exit(0)
	}

	// binding queue
	for _, v := range os.Args[1:] {
		log.Printf("Binding queue[%s] to exchange[%s] with routing key[%s]\n", queue.Name, xName, v)

		err = ch.QueueBind(queue.Name, v, xName, false, nil)
		utils.FatalOnError(err, "failed to bind a queue")
	}

	// delivery msg
	delivery, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	utils.FatalOnError(err, "failed to register a consume")

	// read msg
	forever := make(chan struct{})

	go func() {
		for d := range delivery {
			log.Printf("[x] %s", d.Body)
		}
	}()

	log.Printf("[*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
