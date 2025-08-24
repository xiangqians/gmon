// @author xiangqian
// @date 2025/07/26 10:55
package tmpl

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"strings"
)

// 是否是开发环境
const dev = false

//go:embed html/*
var embedfs embed.FS

var tmpl *template.Template

var funcMap = template.FuncMap{
	"contains":  strings.Contains,
	"hasPrefix": strings.HasPrefix,
	"hasSuffix": strings.HasSuffix,
}

// Init 初始化模板
func Init() error {
	// 解析模板
	var err error
	tmpl, err = template.New("").Funcs(funcMap).ParseFS(embedfs, "html/*")
	if err != nil {
		return err
	}
	return nil
}

func Execute(w io.Writer, name string, data any) {
	if dev {
		// 从文件系统加载，支持热重载
		tmpl, err := template.New("").Funcs(funcMap).ParseGlob("pkg/tmpl/html/*")
		if err != nil {
			log.Println(err)
			return
		}

		// 执行模板
		err = tmpl.ExecuteTemplate(w, fmt.Sprintf("%s.html", name), data)
		if err != nil {
			log.Println(err)
		}
		return
	}

	// 执行模板
	err := tmpl.ExecuteTemplate(w, fmt.Sprintf("%s.html", name), data)
	if err != nil {
		log.Println(err)
	}
}
