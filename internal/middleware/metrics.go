package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"maps"
	"os"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"text/tabwriter"
	"time"
)

type MetricsData struct {
	TotalRequests      uint64
	TotalBytesReceived uint64
	TotalBytesSent     uint64
	Total1xx           uint64
	Total2xx           uint64
	Total3xx           uint64
	Total4xx           uint64
	Total5xx           uint64
	TotalDuration      uint64
}

type MetricsStore struct {
	metricsMap sync.Map // map[string]*MetricsData
}

func (store *MetricsStore) AddMetric(url, method string, recd, sent int64, duration time.Duration, status int) {
	key := url + "|" + method
	metrics, _ := store.metricsMap.LoadOrStore(key, &MetricsData{})
	md := metrics.(*MetricsData)

	atomic.AddUint64(&md.TotalRequests, 1)
	atomic.AddUint64(&md.TotalDuration, uint64(duration.Nanoseconds()))
	if recd < 0 {
		recd = 0
	}
	atomic.AddUint64(&md.TotalBytesReceived, uint64(recd))
	if sent < 0 {
		sent = 0
	}
	atomic.AddUint64(&md.TotalBytesSent, uint64(sent))

	switch {
	case status >= 100 && status < 200:
		atomic.AddUint64(&md.Total1xx, 1)
	case status >= 200 && status < 300:
		atomic.AddUint64(&md.Total2xx, 1)
	case status >= 300 && status < 400:
		atomic.AddUint64(&md.Total3xx, 1)
	case status >= 400 && status < 500:
		atomic.AddUint64(&md.Total4xx, 1)
	case status >= 500:
		atomic.AddUint64(&md.Total5xx, 1)
	}
}

func (store *MetricsStore) GetMetrics() map[string]MetricsData {
	result := make(map[string]MetricsData)
	store.metricsMap.Range(func(key, data interface{}) bool {
		result[key.(string)] = MetricsData{
			TotalRequests:      atomic.LoadUint64(&data.(*MetricsData).TotalRequests),
			TotalBytesReceived: atomic.LoadUint64(&data.(*MetricsData).TotalBytesReceived),
			TotalBytesSent:     atomic.LoadUint64(&data.(*MetricsData).TotalBytesSent),
			Total1xx:           atomic.LoadUint64(&data.(*MetricsData).Total1xx),
			Total2xx:           atomic.LoadUint64(&data.(*MetricsData).Total2xx),
			Total3xx:           atomic.LoadUint64(&data.(*MetricsData).Total3xx),
			Total4xx:           atomic.LoadUint64(&data.(*MetricsData).Total4xx),
			Total5xx:           atomic.LoadUint64(&data.(*MetricsData).Total5xx),
			TotalDuration:      atomic.LoadUint64(&data.(*MetricsData).TotalDuration),
		}
		return true
	})
	return result
}

func (store *MetricsStore) WriteMetrics() {
	fmt.Print("\033[H\033[2J")
	headers := []string{
		"URL",
		"METHOD",
		"TOTAL REQUESTS",
		"1xx",
		"2xx",
		"3xx",
		"4xx",
		"5xx",
		"BYTES RECD",
		"BYTES SENT",
		"AVG BYTES RECD",
		"AVG BYTES SENT",
		"AVG DURATION",
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	metrics := store.GetMetrics()
	keys := maps.Keys(metrics)
	fmt.Fprintln(w, strings.Join(headers, "\t"))
	for _, key := range slices.Sorted(keys) {
		data := metrics[key]
		row := []string{
			strings.Split(key, "|")[0],
			strings.Split(key, "|")[1],
			fmt.Sprintf("%d", data.TotalRequests),
			fmt.Sprintf("%d", data.Total1xx),
			fmt.Sprintf("%d", data.Total2xx),
			fmt.Sprintf("%d", data.Total3xx),
			fmt.Sprintf("%d", data.Total4xx),
			fmt.Sprintf("%d", data.Total5xx),
			byteSI(data.TotalBytesReceived),
			byteSI(data.TotalBytesSent),
			byteSI(data.TotalBytesReceived / data.TotalRequests),
			byteSI(data.TotalBytesSent / data.TotalRequests),
			data.duration(),
		}
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
	w.Flush()
}

func (data *MetricsData) duration() string {
	totalRequests := atomic.LoadUint64(&data.TotalRequests)
	if totalRequests == 0 {
		return "0s"
	}
	totalDuration := atomic.LoadUint64(&data.TotalDuration)
	avgDuration := time.Duration(totalDuration / totalRequests)
	return avgDuration.String()
}

func byteSI(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

var MetricsResults = &MetricsStore{}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		MetricsResults.AddMetric(c.Request.URL.Path, c.Request.Method, c.Request.ContentLength, int64(c.Writer.Size()), time.Since(startTime), c.Writer.Status())
	}
}
