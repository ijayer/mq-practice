/*
 * 说明：
 * 作者：jayer
 * 时间：2019-06-22 3:59 PM
 * 更新：
 */

package main

import (
	"log"
	"os"

	"github.com/ijayer/mq-practice/utils"
	"github.com/streadway/amqp"
)

func main() {
	// create a connection with rabbitmq
	log.Printf("connecting to rabbitmq: %v\n", utils.Host)
	conn, err := amqp.Dial(utils.Host)
	utils.FatalOnError(err, "failed to create a connection with rabbitmq")
	log.Printf("connection established")
	defer conn.Close()

	// create a communication channel
	ch, err := conn.Channel()
	utils.FatalOnError(err, "failed to create a channel")
	defer ch.Close()

	// create an exchange
	xName := "topic_logs"
	xKind := "topic"
	err = ch.ExchangeDeclare(xName, xKind, true, false, false, false, nil)
	utils.FatalOnError(err, "failed to create an exchange")

	// create a queue
	queue, err := ch.QueueDeclare("logs", true, false, false, false, nil)
	utils.FatalOnError(err, "failed to declare a queue")

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [binding_key]...", os.Args[0])
		os.Exit(0)
	}

	// binding queue with exchange
	for _, s := range os.Args[1:] {
		log.Printf("Binding queue[%s] to exchange[%s] with routing_key[%s]\n", queue.Name, xName, s)
		err = ch.QueueBind(queue.Name, s, xName, false, nil)
		utils.FatalOnError(err, "failed to bind a queue")
	}

	// register a consumer to the queue
	delivery, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	utils.FatalOnError(err, "failed to register a consumer to the queue")

	// starting consume msg on a goroutine
	go func() {
		for msg := range delivery {
			log.Printf("[x] msg: %v\n", string(msg.Body))
		}
	}()
	log.Printf("[*] Waiting for logs. To exit press CTRL+C")

	for {
		select {}
	}
}

/*
 How to recv msg?

	To receive all the logs:
		$ go run recv_logs_topic.go "#"
	To receive all logs from the facility "kern":
		$ go run recv_log_topic.go "kern.*"
	To only receive critical logs:
		$ go run recv_log_topic.go "*.critical"
	To receive logs with multipart bindings:
		$ go run recv_log_topic.go "kern.*" "*.critical"
*/
