// @author xiangqian
// @date 2025/07/26 19:31
package handler

import (
	"fmt"
	"gmon/pkg/prom"
	"gmon/pkg/tmpl"
	"net/http"
)

func index(prefix, user string, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "" || r.URL.Path == fmt.Sprintf("/%s", prefix) {
		var data = make(map[string]any)
		data["prefix"] = prefix
		data["user"] = user

		apps, err := prom.Apps()
		if err != nil {
			data["error"] = err.Error()
		}
		data["apps"] = apps

		tmpl.Execute(w, "index", data)
		return
	}

	tmpl.Execute(w, "error", map[string]any{
		"prefix": prefix,
		"error":  "NotFound",
	})
}
