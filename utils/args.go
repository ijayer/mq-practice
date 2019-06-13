/*
 * 说明：
 * 作者：zhe
 * 时间：2019-03-29 4:18 PM
 * 更新：
 */

package utils

import (
	"os"
	"strings"
)

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

// SeverityForm 从命令行输入获取日志级别
func SeverityForm(args []string) string {
	var s string
	var l = len(args)
	if l < 2 || os.Args[1] == "" {
		s = "info"
	} else {
		s = os.Args[1]
	}
	return s
}
