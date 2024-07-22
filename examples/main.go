package main

import (
	"github.com/168yy/gf-metrics"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"time"
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
	r.POST("/slow5s", func(r *ghttp.Request) {
		time.Sleep(time.Duration(5) * time.Second)
		r.Response.Write("Sleep 5s ok")
	})
	r.GET("/slow10s", func(r *ghttp.Request) {
		time.Sleep(time.Duration(10) * time.Second)
		r.Response.Write("Sleep 10s ok")
	})
	r.POST("/slow20s", func(r *ghttp.Request) {
		time.Sleep(time.Duration(20) * time.Second)
		r.Response.Write("Sleep 20s ok")
	})
	r.GET("/slow30s", func(r *ghttp.Request) {
		time.Sleep(time.Duration(30) * time.Second)
		r.Response.Write("Sleep 30s ok")
	})
	r.POST("/slow60s", func(r *ghttp.Request) {
		time.Sleep(time.Duration(60) * time.Second)
		r.Response.Write("Sleep 60s ok")
	})
	s.Run()
}
