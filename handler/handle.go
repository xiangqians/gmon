// @author xiangqian
// @date 2025/07/26 11:49
package handler

import (
	"fmt"
	"gmon/pkg/xhttp"
	"net/http"
)

func Handle(prefix, user, passwd string) {
	xhttp.Handle(prefix, "/", func(w http.ResponseWriter, r *http.Request) {
		index(prefix, user, w, r)
	})
	http.HandleFunc(fmt.Sprintf("%s/login", prefix), func(w http.ResponseWriter, r *http.Request) {
		login(prefix, w, r)
	})
	http.HandleFunc(fmt.Sprintf("%s/login1", prefix), func(w http.ResponseWriter, r *http.Request) {
		login1(prefix, user, passwd, w, r)
	})
	xhttp.Handle(prefix, "/logout", func(w http.ResponseWriter, r *http.Request) {
		logout(prefix, w, r)
	})
	xhttp.Handle(prefix, "/event", event)
}
