/*
 * 说明：
 * 作者：jayer
 * 时间：2019-06-22 3:58 PM
 * 更新：
 */

package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/ijayer/mq-practice/utils"
	"github.com/streadway/amqp"
)

type parser struct{}

func (p parser) BodyFrom(args []string) string {
	var s string
	if len(os.Args) < 3 || os.Args[2] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[2:], " ")
	}

	return s
}

func (p parser) SeverityFrom(args []string) string {
	var s string
	if len(os.Args) < 2 || os.Args[1] == "" {
		s = "anonymous.info"
	} else {
		s = os.Args[1]
	}
	return s
}

func main() {
	// create a conn with rabbitmq
	log.Printf("connecting to rabbitmq[%v]\n", utils.Host)
	conn, err := amqp.Dial(utils.Host)
	utils.FatalOnError(err, "failed to connect with rabbitmq:"+utils.Host)
	defer conn.Close()
	log.Printf("connection established")

	// create the channel for communication
	ch, err := conn.Channel()
	utils.FatalOnError(err, "failed to create channel")
	defer ch.Close()

	// define an exchange
	xName := "topic_logs"
	xKind := "topic"
	err = ch.ExchangeDeclare(xName, xKind, true, false, false, false, nil)
	utils.FatalOnError(err, "failed to define an exchange")

	p := parser{}
	routingKey := p.SeverityFrom(os.Args)
	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(p.BodyFrom(os.Args)),
	}

	start := time.Now()
	err = ch.Publish(xName, routingKey, false, false, msg)
	utils.FatalOnError(err, "failed to publish a msg")
	duration := time.Now().Sub(start)
	log.Printf("[T] time cost for publishing msg: %v\n", duration.String())

	log.Printf("[X] sent msg: %v\n", p.BodyFrom(os.Args))
}

/*
	How to emit logs?

		emit a log with routing key: "kern.critical"
			$ go run emit_log_topic.go "kern.critical"
*/
