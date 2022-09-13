package metrics

import (
	"fmt"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/jxo-me/gf-metrics/bloom"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricType int

const (
	None MetricType = iota
	Counter
	Gauge
	Histogram
	Summary

	DefaultMetricPath = "/metrics"
	DefaultSlowTime   = int32(5)
)

var (
	defaultDuration = []float64{0.1, 0.3, 1.2, 5, 10}
	insMonitor      = Monitor{
		metricPath:  DefaultMetricPath,
		slowTime:    DefaultSlowTime,
		reqDuration: defaultDuration,
		metrics:     make(map[string]*Metric),
	}

	promTypeHandler = map[MetricType]func(metric *Metric) error{
		Counter:   counterHandler,
		Gauge:     gaugeHandler,
		Histogram: histogramHandler,
		Summary:   summaryHandler,
	}
)

// Monitor is an object that uses to set Gf server monitor.
type Monitor struct {
	slowTime    int32
	metricPath  string
	reqDuration []float64
	metrics     map[string]*Metric
}

// GetMonitor used to get global Monitor object,
// this function returns a singleton object.
func GetMonitor() *Monitor {
	return &insMonitor
}

// GetMetric used to get metric object by metric_name.
func (m *Monitor) GetMetric(name string) *Metric {
	if metric, ok := m.metrics[name]; ok {
		return metric
	}
	return &Metric{}
}

// SetMetricPath set metricPath property. metricPath is used for Prometheus
// to get Gf server monitoring data.
func (m *Monitor) SetMetricPath(path string) {
	m.metricPath = path
}

// SetSlowTime set slowTime property. slowTime is used to determine whether
// the request is slow. For "gf_slow_request_total" metric.
func (m *Monitor) SetSlowTime(slowTime int32) {
	m.slowTime = slowTime
}

// SetDuration set reqDuration property. reqDuration is used to ginRequestDuration
// metric buckets.
func (m *Monitor) SetDuration(duration []float64) {
	m.reqDuration = duration
}

func (m *Monitor) SetMetricPrefix(prefix string) {
	metricRequestTotal = prefix + metricRequestTotal
	metricRequestUVTotal = prefix + metricRequestUVTotal
	metricURIRequestTotal = prefix + metricURIRequestTotal
	metricRequestBody = prefix + metricRequestBody
	metricResponseBody = prefix + metricResponseBody
	metricRequestDuration = prefix + metricRequestDuration
	metricSlowRequest = prefix + metricSlowRequest
}

func (m *Monitor) SetMetricSuffix(suffix string) {
	metricRequestTotal += suffix
	metricRequestUVTotal += suffix
	metricURIRequestTotal += suffix
	metricRequestBody += suffix
	metricResponseBody += suffix
	metricRequestDuration += suffix
	metricSlowRequest += suffix
}

// AddMetric add custom monitor metric.
func (m *Monitor) AddMetric(metric *Metric) error {
	if _, ok := m.metrics[metric.Name]; ok {
		return errors.Errorf("metric '%s' is existed", metric.Name)
	}

	if metric.Name == "" {
		return errors.Errorf("metric name cannot be empty.")
	}
	if f, ok := promTypeHandler[metric.Type]; ok {
		if err := f(metric); err == nil {
			prometheus.MustRegister(metric.vec)
			m.metrics[metric.Name] = metric
			return nil
		}
	}
	return errors.Errorf("metric type '%d' not existed.", metric.Type)
}

// InitMetrics used to init Gf metrics
func (m *Monitor) InitMetrics() *Monitor {
	bloomFilter = bloom.NewBloomFilter()

	_ = m.AddMetric(&Metric{
		Type:        Counter,
		Name:        metricRequestTotal,
		Description: "all the server received request num.",
		Labels:      nil,
	})
	_ = m.AddMetric(&Metric{
		Type:        Counter,
		Name:        metricRequestUVTotal,
		Description: "all the server received ip num.",
		Labels:      nil,
	})
	_ = m.AddMetric(&Metric{
		Type:        Counter,
		Name:        metricURIRequestTotal,
		Description: "all the server received request num with every uri.",
		Labels:      []string{"uri", "method", "code"},
	})
	_ = m.AddMetric(&Metric{
		Type:        Counter,
		Name:        metricRequestBody,
		Description: "the server received request body size, unit byte",
		Labels:      nil,
	})
	_ = m.AddMetric(&Metric{
		Type:        Counter,
		Name:        metricResponseBody,
		Description: "the server send response body size, unit byte",
		Labels:      nil,
	})
	_ = m.AddMetric(&Metric{
		Type:        Histogram,
		Name:        metricRequestDuration,
		Description: "the time server took to handle the request.",
		Labels:      []string{"uri"},
		Buckets:     m.reqDuration,
	})
	_ = m.AddMetric(&Metric{
		Type:        Counter,
		Name:        metricSlowRequest,
		Description: fmt.Sprintf("the server handled slow requests counter, t=%d.", m.slowTime),
		Labels:      []string{"uri", "method", "code"},
	})

	return m
}

// Use set Gf metrics middleware
func (m *Monitor) Use(s *ghttp.Server) {
	m.InitMetrics()
	s.BindMiddlewareDefault(m.Middleware)
	s.BindHandler(m.metricPath, ghttp.WrapH(promhttp.Handler()))
}

// UseWithoutExposingEndpoint is used to add monitor interceptor to Gf router
// It can be called multiple times to intercept from multiple gin.IRoutes
// http path is not set, to do that use Expose function
func (m *Monitor) UseWithoutExposingEndpoint(s *ghttp.Server) {
	m.InitMetrics()
	s.BindMiddlewareDefault(m.Middleware)
}
