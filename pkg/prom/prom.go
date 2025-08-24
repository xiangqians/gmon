// @author xiangqian
// @date 2025/07/26 13:11
package prom

import (
	"context"
	"fmt"
	pkg_api "github.com/prometheus/client_golang/api"
	pkg_api_v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"gmon/pkg/xtime"
	"log"
	"math"
	"sort"
	"time"
)

// Prometheus 普罗米修斯

// Exporter
// https://prometheus.io/docs/instrumenting/exporters/

// 通过 Prometheus Web UI 查看 Prometheus 默认数据保留时间
// 访问 Prometheus 的 Web 界面（默认 http://<prom-server>:9090）
// 导航到 Status -> Runtime & Build Information
// 查找 Storage retention	15d

var api pkg_api_v1.API

func Init(config Config) error {
	// 创建 Prometheus 客户端
	client, err := pkg_api.NewClient(pkg_api.Config{
		Address: fmt.Sprintf("http://%s:%d", config.Host, config.Port),
	})
	if err != nil {
		return err
	}

	api = pkg_api_v1.NewAPI(client)

	// 查询 Prometheus 自身的状态指标
	ctx, cancel := withTimeout()
	defer cancel()
	_, _, err = api.Query(ctx, "up", time.Now())
	if err != nil {
		return err
	}
	return nil
}

func Apps() ([]*App, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	// 查询所有服务器
	targets, err := api.Targets(ctx)
	if err != nil {
		return nil, err
	}

	active := targets.Active
	apps := make([]*App, 0, len(active))
label:
	for _, act := range active {
		var appName = string(act.Labels["app"])
		var instName = string(act.Labels["job"])
		var instAddr = string(act.Labels["instance"])

		var status Status
		var tm time.Time
		var duration time.Duration
		switch act.Health {
		case pkg_api_v1.HealthGood:
			status = StatusUp
			start, _ := LastDownTime(instName, instAddr)
			if start.IsZero() {
				start, _ = FirstUpTime(instName, instAddr)
			}
			tm = start
			end, _ := LastUpTime(instName, instAddr)
			duration = end.Sub(start)
			if duration < 0 {
				duration = 0
			}

		case pkg_api_v1.HealthBad:
			status = StatusDown
			tm, _ = LastUpTime(instName, instAddr)
			if tm.IsZero() {
				tm, _ = FirstDownTime(instName, instAddr)
			}
		}

		var instance = &Instance{
			Name:     instName,
			Addr:     instAddr,
			Status:   status,
			Time:     xtime.XTime{Time: tm},
			Duration: xtime.XDuration{Duration: duration},
		}

		for _, app := range apps {
			if app.Name == appName {
				app.Instances = append(app.Instances, instance)
				continue label
			}
		}

		app := &App{
			Name:      appName,
			Instances: []*Instance{instance},
		}
		apps = append(apps, app)
	}

	// 排序
	sort.Slice(apps, func(i, j int) bool {
		return number(apps[i].Instances[0].Name) < number(apps[j].Instances[0].Name)
	})
	for _, app := range apps {
		sort.Slice(app.Instances, func(i, j int) bool {
			return app.Instances[i].Addr < app.Instances[j].Addr
		})
	}

	return apps, nil
}

func number(name string) uint8 {
	switch name {
	case "go":
		return 1
	case "java":
		return 2
	case "mysql":
		return 3
	case "redis":
		return 4
	case "windows":
		return 5
	case "linux":
		return 6
	case "prom":
		return math.MaxUint8
	default:
		return math.MaxUint8 - 1
	}
}

// FirstUpTime 应用实例最早一次在线时间
func FirstUpTime(name, addr string) (time.Time, error) {
	vector, err := Vector(fmt.Sprintf(`min_over_time(timestamp(up{job="%s", instance="%s"} == 1)[15d:])`, name, addr))
	if err != nil {
		return time.Time{}, nil
	}

	if vector.Len() == 0 {
		return time.Time{}, nil
	}

	timestamp := int64(vector[0].Value)
	return time.Unix(timestamp, 0), nil
}

// LastUpTime 应用实例最近一次在线时间
func LastUpTime(name, addr string) (time.Time, error) {
	vector, err := Vector(fmt.Sprintf(`max_over_time(timestamp(up{job="%s", instance="%s"} == 1)[15d:])`, name, addr))
	if err != nil {
		return time.Time{}, nil
	}

	if vector.Len() == 0 {
		return time.Time{}, nil
	}

	timestamp := int64(vector[0].Value)
	return time.Unix(timestamp, 0), nil
}

// FirstDownTime 应用实例最早一次离线时间
func FirstDownTime(name, addr string) (time.Time, error) {
	vector, err := Vector(fmt.Sprintf(`min_over_time(timestamp(up{job="%s", instance="%s"} == 0)[15d:])`, name, addr))
	if err != nil {
		return time.Time{}, nil
	}

	if vector.Len() == 0 {
		return time.Time{}, nil
	}

	timestamp := int64(vector[0].Value)
	return time.Unix(timestamp, 0), nil
}

// LastDownTime 应用实例最近一次离线时间
func LastDownTime(name, addr string) (time.Time, error) {
	vector, err := Vector(fmt.Sprintf(`max_over_time(timestamp(up{job="%s", instance="%s"} == 0)[15d:])`, name, addr))
	if err != nil {
		return time.Time{}, nil
	}

	if vector.Len() == 0 {
		return time.Time{}, nil
	}

	timestamp := int64(vector[0].Value)
	return time.Unix(timestamp, 0), nil
}

// LastSample 最新采样
func LastSample(expr string) (*Sample, error) {
	log.Printf("expr: %s\n", expr)

	vector, err := Vector(expr)
	if err != nil {
		return nil, err
	}

	var sample *Sample = nil
	for _, samp := range vector {
		metric := samp.Metric
		var instAddr = string(metric["instance"])
		var name = string(metric["name"])
		var key = fmt.Sprintf("%s,%s", instAddr, name)
		var value = float64(samp.Value)
		if sample == nil {
			sample = &Sample{
				Timestamp: int64(samp.Timestamp),
				Value:     make(map[string]float64),
			}
		}
		sample.Value[key] = value
	}
	return sample, nil
}

func Vector(expr string) (model.Vector, error) {
	value, err := Query(expr)
	if err != nil {
		return nil, err
	}

	vector, ok := value.(model.Vector)
	if !ok {
		return nil, fmt.Errorf("cannot convert result to vector")
	}

	return vector, nil
}

func Query(expr string) (model.Value, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	// PromQL 的 [range:offset] 语法：[查询的时间窗口长度, 相对于评估时间点的偏移量]
	// 计算方式：查询时间范围 = [评估时间 - range - offset, 评估时间 - offset]
	value, _, err := api.Query(ctx,
		expr,       // 表达式
		time.Now()) // 评估时间
	return value, err
}

func withTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

// App 应用
type App struct {
	Name      string      `json:"name"`      // 名称
	Instances []*Instance `json:"instances"` // 实例集
}

// Instance 实例
type Instance struct {
	Name     string          `json:"name"`     // 名称
	Addr     string          `json:"addr"`     // 地址
	Status   Status          `json:"status"`   // 状态
	Time     xtime.XTime     `json:"time"`     // 在线/离线时间
	Duration xtime.XDuration `json:"duration"` // 在线持续时间
}

// Sample 采样
type Sample struct {
	Timestamp int64              `json:"timestamp"`
	Value     map[string]float64 `json:"value"`
}

type Status byte

const (
	StatusUp Status = iota + 1
	StatusDown
)

func (status Status) MarshalJSON() ([]byte, error) {
	return []byte(`"` + status.String() + `"`), nil
}

func (status Status) String() string {
	switch status {
	case StatusUp:
		return "UP"
	case StatusDown:
		return "DOWN"
	default:
		return "UNKNOWN"
	}
}

// Config Prometheus 配置
type Config struct {
	Host string // Prometheus 主机
	Port uint16 // Prometheus 端口
}
