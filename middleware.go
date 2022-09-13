package metrics

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"strconv"
	"time"
)

// Middleware as Gf monitor middleware.
func (m *Monitor) Middleware(r *ghttp.Request) {
	if r.Request.URL.Path == m.metricPath {
		r.Middleware.Next()
		return
	}
	startTime := time.Now()

	// execute normal process.
	r.Middleware.Next()

	// after request
	m.MetricHandle(r, startTime)
}

func (m *Monitor) MetricHandle(r *ghttp.Request, start time.Time) {
	// set request total
	_ = m.GetMetric(metricRequestTotal).Inc(nil)

	// set uv
	if clientIP := r.GetClientIp(); !bloomFilter.Contains(clientIP) {
		bloomFilter.Add(clientIP)
		_ = m.GetMetric(metricRequestUVTotal).Inc(nil)
	}

	// set uri request total
	_ = m.GetMetric(metricURIRequestTotal).Inc([]string{r.Request.URL.Path, r.Method, strconv.Itoa(r.Response.Status)})

	// set request body size
	// since r.ContentLength can be negative (in some occasions) guard the operation
	if r.ContentLength >= 0 {
		_ = m.GetMetric(metricRequestBody).Add(nil, float64(r.ContentLength))
	}

	// set slow request
	latency := time.Since(start)
	if int32(latency.Seconds()) > m.slowTime {
		_ = m.GetMetric(metricSlowRequest).Inc([]string{r.Request.URL.Path, r.Method, strconv.Itoa(r.Response.Status)})
	}

	// set request duration
	_ = m.GetMetric(metricRequestDuration).Observe([]string{r.Request.URL.Path}, latency.Seconds())

	// set response size
	if r.Response.BufferLength() > 0 {
		_ = m.GetMetric(metricResponseBody).Add(nil, float64(r.Response.BufferLength()))
	}
	glog.Debug(r.GetCtx(), "Response body len:", r.Response.BufferLength())
}
