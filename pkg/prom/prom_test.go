// @author xiangqian
// @date 2025/07/26 13:12
package prom

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	err := Init(Config{"localhost", 9090})
	if err != nil {
		panic(err)
	}
}

func TestApps(t *testing.T) {
	TestInit(t)
	apps, err := Apps()
	if err != nil {
		panic(err)
	}
	Print(apps)
}

func TestQuery(t *testing.T) {
	TestInit(t)
	expr := `delta(
  label_replace(
    sum by (job, instance) (go_memstats_sys_bytes{job="prom"}),
    "name", "mem_used_bytes", "", ""
  )[1h:]
)`

	expr = `go_memstats_sys_bytes{job="go"}[1m]`

	// go_memstats_sys_bytes{job="prom"}[1h]

	value, err := Query(expr)
	if err != nil {
		panic(err)
	}
	Print(value)
}

func TestLastCpuMem(t *testing.T) {
	TestInit(t)

	var expr []string

	sample, err := LastSample(strings.Join(expr, " or "))
	if err != nil {
		panic(err)
	}
	Print(sample)
}

// 查询系统进程内存
func Test_process_resident_memory_bytes(t *testing.T) {
	TestInit(t)
	var expr = "process_resident_memory_bytes"
	expr = `100 - (avg by (job,instance) (rate(windows_cpu_time_total{mode="idle"}[10s])) * 100)`
	expr += ` or process_resident_memory_bytes`
	sample, err := LastSample(expr)
	if err != nil {
		panic(err)
	}
	Print(sample)
}
func Print(v any) {
	//fmt.Printf("%+v\n", v)
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", data)
}
