// @author xiangqian
// @date 2025/07/20 16:08
package xhttp

import (
	"fmt"
	"net/http"
)

func Handle(prefix, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(fmt.Sprintf("%s%s", prefix, pattern), func(w http.ResponseWriter, r *http.Request) {
		// 会话是否已过期
		if Expired(r) {
			// 重定向到登录页
			http.Redirect(w, r, fmt.Sprintf("%s/login", prefix), http.StatusFound)
			return
		}

		// 会话有效，继续处理请求
		handler(w, r)
	})
}
