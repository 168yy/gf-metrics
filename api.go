package metrics

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ApiExpose adds metric path to a given router.
// The router can be different with the one passed to UseWithoutExposingEndpoint.
// This allows to expose metrics on different port.
func (m *Monitor) ApiExpose(r *ghttp.Request) {
	promhttp.Handler().ServeHTTP(r.Response.ResponseWriter, r.Request)
}
