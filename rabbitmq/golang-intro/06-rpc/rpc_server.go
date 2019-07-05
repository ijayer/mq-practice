/*
 * 说明：
 * 作者：jayer
 * 时间：2019-06-23 11:16 AM
 * 更新：
 */

package main

import (
	"log"
	"strconv"

	"github.com/ijayer/mq-practice/utils"
	"github.com/streadway/amqp"
)

func fib(n int) int {
	if n == 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	return fib(n-1) + fib(n-2)
}

func main() {
	// create connection with rabbitmq
	log.Printf("connecting to rabbitmq[%s]\n", utils.Host)
	conn, err := amqp.Dial(utils.Host)
	utils.FatalOnError(err, "failed to connect")
	defer conn.Close()
	log.Printf("connection established")

	// create communication channel
	ch, err := conn.Channel()
	utils.FatalOnError(err, "failed to create a channel")
	defer ch.Close()

	// queue declare
	queueName := "rpc_queue"
	q, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	utils.FatalOnError(err, "failed to declare a queue")

	err = ch.Qos(1, 0, false)
	utils.FatalOnError(err, "failed to set Qos")

	// consume msg from queue
	delivery, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	utils.FatalOnError(err, "failed to register a consume")

	// handle request and return response in a goroutine
	forever := make(chan bool)

	go func() {
		for msg := range delivery {
			n, err := strconv.Atoi(string(msg.Body))
			utils.FatalOnError(err, "failed to convert body to integer")

			log.Printf("[.] fib(%d)", n)
			response := fib(n)

			// publish resp msg to client
			err = ch.Publish(
				"",
				msg.ReplyTo,
				false,
				false,
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: msg.CorrelationId,
					Body:          []byte(strconv.Itoa(response)),
				},
			)
			utils.FatalOnError(err, "failed to publish a msg")
			log.Printf("[.] response: %d\n", response)

			msg.Ack(false)
		}
	}()

	log.Printf("[*] awaiting RPC requests")

	<-forever
}

/*
    Start the rpc Server:

	$ go run rpc_server.go
*/
