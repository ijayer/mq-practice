/*
 * 说明：
 * 作者：jayer
 * 时间：2019-06-23 11:16 AM
 * 更新：
 */

package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ijayer/mq-practice/utils"
	"github.com/streadway/amqp"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	n := bodyFrom(os.Args)

	log.Printf("[x] Requesting fib(%d)", n)
	resp, err := fibonacciRPC(n)
	utils.FatalOnError(err, "failed to handle rpc request")

	log.Printf("[.] Got %d\n", resp)
}

func bodyFrom(args []string) int {
	var s string
	if len(args) < 2 || os.Args[1] == "" {
		s = "30"
	} else {
		s = strings.Join(args[1:], " ")
	}
	n, err := strconv.Atoi(s)
	utils.FatalOnError(err, "failed to convert args to integer")
	return n
}

func randomStr(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func fibonacciRPC(n int) (res int, err error) {
	// create a connection with rabbitmq
	log.Printf("connecting to rabbitmq[%s]\n", utils.Host)
	conn, err := amqp.Dial(utils.Host)
	utils.FatalOnError(err, "failed to connect")
	defer conn.Close()
	log.Printf("connection established")

	// create a channel
	ch, err := conn.Channel()
	utils.FatalOnError(err, "failed to create a channel")
	defer ch.Close()

	// declare a queue
	queue, err := ch.QueueDeclare("", false, false, true, false, nil)
	utils.FatalOnError(err, "failed to declare a queue")

	// register a consumer
	delivery, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	utils.FatalOnError(err, "failed to register a consumer")

	// publish request
	corrId := randomStr(32)
	err = ch.Publish(
		"",
		"rpc_queue",
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       queue.Name,
			Body:          []byte(strconv.Itoa(n)),
		})

	// recv response
	for msg := range delivery {
		if corrId == msg.CorrelationId {
			res, err = strconv.Atoi(string(msg.Body))
			utils.FatalOnError(err, "failed to convert body to integer")
			break
		}
	}

	return
}

/*
	Start the rcp client: calculate a fibonacci number

	$ go run rpc_client.go 30
*/
