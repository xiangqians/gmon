// @author xiangqian
// @date 2025/07/26 13:57
package xtime

import (
	"fmt"
	"time"
)

func Format(time time.Time) string {
	if time.IsZero() {
		return "--"
	}
	return time.Format("2006/01/02 15:04:05")
}

type XTime struct {
	time.Time
}

func (xtime XTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + xtime.String() + `"`), nil
}

func (xtime XTime) String() string {
	return Format(xtime.Time)
}

type XDuration struct {
	time.Duration
}

func (xduration XDuration) MarshalJSON() ([]byte, error) {
	return []byte(`"` + xduration.String() + `"`), nil
}

func (xduration XDuration) String() string {
	if xduration.Duration <= 0 {
		return "--"
	}

	duration := xduration.Duration.Round(time.Second)

	// 天
	d := duration / (24 * time.Hour)
	if d > 0 {
		h := (duration - d*(24*time.Hour)) / time.Hour
		if h > 0 {
			return fmt.Sprintf("%dd%dh", d, h)
		}
		return fmt.Sprintf("%dd", d)
	}

	// 小时
	h := duration / time.Hour
	if h > 0 {
		m := (duration - h*time.Hour) / time.Minute
		if m > 0 {
			return fmt.Sprintf("%dh%dm", h, m)
		}
		return fmt.Sprintf("%dh", h)
	}

	// 分钟
	m := duration / time.Minute
	if m > 0 {
		return fmt.Sprintf("%dm", m)
	}

	return xduration.Duration.String()
}
