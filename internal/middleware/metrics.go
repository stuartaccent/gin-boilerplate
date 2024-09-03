package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"sort"
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
	data sync.Map // map[string]*MetricsData
}

func (store *MetricsStore) AddMetric(url, method string, recd, sent int64, duration time.Duration, status int) {
	key := url + "|" + method
	metrics, _ := store.data.LoadOrStore(key, &MetricsData{})
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

func (store *MetricsStore) SortedMetrics() []struct {
	Key  string
	Data MetricsData
} {
	var metricsSlice []struct {
		Key  string
		Data MetricsData
	}
	store.data.Range(func(key, data interface{}) bool {
		md := data.(*MetricsData)
		metricsSlice = append(metricsSlice, struct {
			Key  string
			Data MetricsData
		}{
			Key: key.(string),
			Data: MetricsData{
				TotalRequests:      atomic.LoadUint64(&md.TotalRequests),
				TotalBytesReceived: atomic.LoadUint64(&md.TotalBytesReceived),
				TotalBytesSent:     atomic.LoadUint64(&md.TotalBytesSent),
				Total1xx:           atomic.LoadUint64(&md.Total1xx),
				Total2xx:           atomic.LoadUint64(&md.Total2xx),
				Total3xx:           atomic.LoadUint64(&md.Total3xx),
				Total4xx:           atomic.LoadUint64(&md.Total4xx),
				Total5xx:           atomic.LoadUint64(&md.Total5xx),
				TotalDuration:      atomic.LoadUint64(&md.TotalDuration),
			},
		})
		return true
	})

	sort.Slice(metricsSlice, func(i, j int) bool {
		return metricsSlice[j].Data.TotalRequests < metricsSlice[i].Data.TotalRequests
	})

	return metricsSlice
}

func (store *MetricsStore) WriteMetrics(topX int) {
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

	fmt.Fprintln(w, strings.Join(headers, "\t"))

	for i, entry := range store.SortedMetrics() {
		if i >= topX {
			break
		}
		data := entry.Data
		row := []string{
			strings.Split(entry.Key, "|")[0],
			strings.Split(entry.Key, "|")[1],
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
	return time.Duration(totalDuration / totalRequests).String()
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
