package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/jxo-me/gf-metrics"
)

func main() {
	s := g.Server()
	s.SetAddr(":8080")
	// get global Monitor object
	m := metrics.GetMonitor()
	// +optional set metric path, default /metrics
	m.SetMetricPath("/metrics")
	// +optional set slow time, default 5s
	m.SetSlowTime(10)
	// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
	// used to p95, p99
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	// set middleware for gf
	m.Use(s)
	r := s.Group("")
	r.GET("/hello", func(r *ghttp.Request) { r.Response.Write("Hello") })
	s.Run()
}
