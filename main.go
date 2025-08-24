// @author xiangqian
// @date 2025/07/26 10:51
package main

import (
	"fmt"
	"gmon/handler"
	"gmon/pkg/prom"
	"gmon/pkg/static"
	"gmon/pkg/tmpl"
	"gmon/pkg/xlog"
	"log"
	"net/http"
)

func main() {
	// [config]
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("load config: %v\n", err)
	}

	// [log]
	xlog.Init()

	// [prom]
	err = prom.Init(config.Prom)
	if err != nil {
		log.Fatalf("init prom: %v\n", err)
	}

	// [static]
	err = static.Init(config.Http.Prefix)
	if err != nil {
		log.Fatalf("init static: %v\n", err)
	}

	// [tmpl]
	err = tmpl.Init()
	if err != nil {
		log.Fatalf("init tmpl: %v\n", err)
	}

	// [handler]
	handler.Handle(config.Http.Prefix, config.Http.User, config.Http.Passwd)

	// 启动服务器
	var port = config.Http.Port
	log.Printf("Server starting on port %d ...\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %v\n", err)
	}
}
