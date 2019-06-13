/*
 * 说明：日志生产者
 * 作者：jayer
 * 时间：2019-06-02 10:53 PM
 * 更新：
 */

package main

import (
	"flag"
	"fmt"
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
	utils.FatalOnError(err, fmt.Sprintf("failed to dial with rabbitmq[%v]", *url))
	defer conn.Close()
	log.Printf("connected")

	// create channel
	ch, err := conn.Channel()
	utils.FatalOnError(err, "failed to create channel")
	defer ch.Close()

	// declare exchange
	xType := "direct"
	xName := "logs_direct"
	err = ch.ExchangeDeclare(xName, xType, true, false, false, false, nil)
	utils.FatalOnError(err, "failed to declare an exchange with type(direct)")

	body := utils.BodyFrom(os.Args)

	// publish msg
	err = ch.Publish(
		xName,
		utils.SeverityForm(os.Args),
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	utils.FatalOnError(err, "failed to publish msg")

	log.Printf("[X] sent: %v\n", body)
}
