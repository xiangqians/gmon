// @author xiangqian
// @date 2025/07/27 16:18
package handler

import (
	"fmt"
	"gmon/pkg/tmpl"
	"gmon/pkg/xhttp"
	"net/http"
)

func login(prefix string, w http.ResponseWriter, r *http.Request) {
	// 会话是否已过期
	if !xhttp.Expired(r) {
		// 没有过期则重定向到首页
		http.Redirect(w, r, fmt.Sprintf("%s/", prefix), http.StatusFound)
		return
	}

	user, _ := xhttp.GetCookie(r, "user")
	err, _ := xhttp.GetCookie(r, "error")
	var data = map[string]any{
		"prefix": prefix,
		"user":   user,
		"error":  err,
	}
	tmpl.Execute(w, "login", data)
	return
}

func login1(prefix, user, passwd string, w http.ResponseWriter, r *http.Request) {
	// 解析表单数据
	err := r.ParseForm()
	if err != nil {
		xhttp.SetCookie(w, "error", err.Error(), 2)
		http.Redirect(w, r, fmt.Sprintf("%s/login", prefix), http.StatusFound)
		return
	}

	ruser := r.FormValue("user")
	rpasswd := r.FormValue("passwd")
	if ruser == user && rpasswd == passwd {
		xhttp.SetSession(w)
		http.Redirect(w, r, fmt.Sprintf("%s/", prefix), http.StatusFound)
		return
	}

	xhttp.SetCookie(w, "user", ruser, 2)
	xhttp.SetCookie(w, "error", "用户名或密码错误", 2)
	http.Redirect(w, r, fmt.Sprintf("%s/login", prefix), http.StatusFound)
}

func logout(prefix string, w http.ResponseWriter, r *http.Request) {
	xhttp.DelSession(w, r)
	http.Redirect(w, r, fmt.Sprintf("%s/login", prefix), http.StatusFound)
}
