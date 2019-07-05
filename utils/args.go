/*
 * 说明：
 * 作者：zhe
 * 时间：2019-03-29 4:18 PM
 * 更新：
 */

package utils

import (
	"flag"
	"os"
	"strings"
)

// rabbitmq server addr
var Host string

func init() {
	flag.StringVar(&Host, "url", "amqp://guest:guest@192.168.0.104:5672/", "rabbitmq server address")
	flag.Parse()
}

type Parser interface {
	BodyFrom([]string) string
	SeverityFrom([]string) string
}

// BodyFrom 获取命令行参数
func BodyFrom(args []string) string {
	var s string
	var l = len(args)

	// os.Args[1] 第一个参数
	if l < 2 || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(os.Args[1:], " ")
	}
	return s
}

// SeverityFrom 从命令行输入获取日志级别
func SeverityFrom(args []string) string {
	var s string
	var l = len(args)
	if l < 2 || os.Args[1] == "" {
		s = "info"
	} else {
		s = os.Args[1]
	}
	return s
}
