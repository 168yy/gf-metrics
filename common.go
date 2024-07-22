package metrics

import "github.com/168yy/gf-metrics/bloom"

var (
	metricRequestTotal    = "gf_request_total"
	metricRequestUVTotal  = "gf_request_uv_total"
	metricURIRequestTotal = "gf_uri_request_total"
	metricRequestBody     = "gf_request_body_total"
	metricResponseBody    = "gf_response_body_total"
	metricRequestDuration = "gf_request_duration"
	metricSlowRequest     = "gf_slow_request_total"

	bloomFilter *bloom.BloomFilter
)
