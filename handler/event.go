// @author xiangqian
// @date 2025/07/27 13:30
package handler

import (
	"fmt"
	"gmon/pkg/prom"
	"gmon/pkg/xjson"
	"net/http"
	"strings"
	"time"
)

//// 常用 Go 内存指标：
////go_memstats_alloc_bytes - 当前分配的堆内存
////go_memstats_sys_bytes - 从系统获取的总内存
////go_memstats_heap_alloc_bytes - 堆内存分配
////process_resident_memory_bytes - 进程实际使用的物理内存
//
////go_memstats_alloc_bytes{job="gweb", instance="localhost:58082"}
//// go_memstats_alloc_bytes{instance="localhost:58082"} or redis_memory_used_bytes{instance="your_redis_host:port"}

////指标名称	说明
////mysql_global_status_innodb_buffer_pool_bytes_total	InnoDB 缓冲池总大小
////mysql_global_status_innodb_buffer_pool_bytes_data	缓冲池中数据占用内存
////mysql_global_status_key_buffer_bytes	MyISAM 键缓存大小
////mysql_global_status_query_cache_size	查询缓存大小
////mysql_global_status_thread_buffers_bytes	线程缓冲区内存
////mysql_global_status_sort_buffer_bytes	排序缓冲区内存

////常用 Redis 内存指标：
////redis_memory_used_bytes - Redis 当前使用的内存量
////redis_memory_max_bytes - 配置的最大内存限制
////redis_memory_peak_bytes - Redis 内存使用的峰值
////# Redis 总内存使用量
////redis_memory_used_bytes{instance="your_redis_host:port"}

func event(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Set("Content-Type", "text/event-stream")
	header.Set("Cache-Control", "no-cache")
	header.Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// 创建退出通道
	done := r.Context().Done()
	for {
		select {
		case <-done:
			return
		default:
		}

		data, err := dat()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data, err = xjson.Serialize(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err = fmt.Fprint(w, fmt.Sprintf("data: %s\n\n", data)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		flusher.Flush()

		time.Sleep(2 * time.Second)
	}
}

func dat() (any, error) {
	apps, err := prom.Apps()
	if err != nil {
		return nil, err
	}

	var expr []string

	// Prometheus
	expr = append(expr, `label_replace(sum by (job, instance) (go_memstats_sys_bytes{job="prom"}), "name", "mem_used_bytes", "", "")`)

	// Windows
	// Windows 系统整体 CPU 使用率（0~100，单位：%）
	expr = append(expr, `label_replace(100 - (avg by (job,instance) (rate(windows_cpu_time_total{job="windows", mode="idle"}[10s])) * 100), "name", "cpu_usage", "", "")`)
	// Windows 系统已使用内存字节数：已使用内存 = 总物理内存 - 可用内存
	expr = append(expr, `label_replace(windows_memory_physical_total_bytes{job="windows"} - windows_memory_physical_free_bytes{job="windows"}, "name", "mem_used_bytes", "", "")`)
	// Windows 内存使用百分比（0~100，单位：%）
	expr = append(expr, `label_replace(((windows_memory_physical_total_bytes{job="windows"} - windows_memory_physical_free_bytes{job="windows"}) / windows_memory_physical_total_bytes{job="windows"}) * 100, "name", "mem_used_percent", "", "")`)

	//总内存和可用内存
	//node_memory_MemTotal_bytes  # 系统总内存
	//node_memory_MemAvailable_bytes  # 可用内存（包含缓存和缓冲）

	//内存使用占比
	//# 已用内存占比（基于 MemAvailable）
	//1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)
	//# 或传统计算方式（used = total - free - buffers - cache）
	//(node_memory_MemTotal_bytes - node_memory_MemFree_bytes - node_memory_Buffers_bytes - node_memory_Cached_bytes) / node_memory_MemTotal_bytes

	//总体 CPU 使用率
	//# 1分钟平均CPU使用率（所有核心）
	//100 - (avg by(instance) (rate(node_cpu_seconds_total{mode="idle"}[1m])) * 100
	//# 5分钟平均CPU使用率
	//100 - (avg by(instance) (rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100

	// Go
	// go_memstats_sys_bytes：Go 总管理内存，Go 从 OS 申请的总内存（含预留）
	// go_memstats_alloc_bytes：Go 实际使用的堆内存，当前存活对象占用的堆内存（不含空闲内存）
	expr = append(expr, `label_replace(sum by (job, instance) (go_memstats_sys_bytes{job="go"}), "name", "mem_used_bytes", "", "")`)

	// Java 已使用内存字节数
	expr = append(expr, `label_replace(sum by (job, instance) (jvm_memory_used_bytes{job="java"}), "name", "mem_used_bytes", "", "")`)

	// MySQL
	// process_resident_memory_bytes RSS（常驻内存）
	expr = append(expr, `label_replace(sum by (job, instance) (process_resident_memory_bytes{job="mysql"}), "name", "mem_used_bytes", "", "")`)

	// Redis
	// 查询 Redis 总内存使用量
	expr = append(expr, `label_replace(sum by (job, instance) (redis_memory_used_bytes{job="redis"}), "name", "mem_used_bytes", "", "")`)

	sample, err := prom.LastSample(strings.Join(expr, " or "))
	if err != nil {
		return nil, err
	}

	return map[string]any{"apps": apps, "sample": sample}, nil
}
