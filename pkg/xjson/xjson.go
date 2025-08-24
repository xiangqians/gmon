// @author xiangqian
// @date 2025/07/20 13:05
package xjson

import (
	"encoding/json"
)

// Serialize JSON 序列化
func Serialize(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Deserialize JSON 反序列化
func Deserialize(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
