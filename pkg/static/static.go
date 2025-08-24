// @author xiangqian
// @date 2025/07/26 11:18
package static

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
)

// 是否是开发环境
const dev = false

//go:embed image/* css/* js/*
var embedfs embed.FS

// Init 初始化静态文件处理器
func Init(prefix string) error {
	for _, dir := range []string{"image", "css", "js"} {
		var iofs fs.FS
		if dev {
			// 从文件系统加载，支持热重载
			iofs = os.DirFS(fmt.Sprintf("pkg/static/%s", dir))
		} else {
			// 使用 fs.Sub 获取子文件系统
			var err error
			iofs, err = fs.Sub(embedfs, dir)
			if err != nil {
				return err
			}
		}

		// 文件服务器
		handler := http.FileServer(http.FS(iofs))
		var pattern = fmt.Sprintf("%s/%s/", prefix, dir)
		http.Handle(pattern, http.StripPrefix(pattern, handler))
	}
	return nil
}
