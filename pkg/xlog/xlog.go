// @author xiangqian
// @date 2025/07/12 10:36
package xlog

import (
	"log"
	"os"
	"time"
)

// Init 初始化日志记录器
func Init() {
	// 将日志输出重定向到 stdout
	//log.SetOutput(os.Stdout)

	// 禁用默认的日期和时间前缀
	log.SetFlags(0)
	// 自定义日志输出
	log.SetOutput(&Logger{})
}

type Logger struct{}

func (logger *Logger) Write(p []byte) (n int, err error) {
	return os.Stdout.WriteString(time.Now().Format("2006/01/02 15:04:05.000") + " " + string(p))
}
