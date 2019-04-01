/*
 * 说明：
 * 作者：zhe
 * 时间：2019-03-28 10:30 PM
 * 更新：
 */

package utils

import (
	"log"
)

// FatalOnError check return error
func FatalOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", err, msg)
	}
}
