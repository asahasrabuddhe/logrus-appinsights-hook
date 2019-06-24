package logrusApplicationInsightsHook

import (
	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
	"time"
)

type Request struct {
	method       string
	uri          string
	duration     time.Duration
	responseCode string
	properties   map[string]string
	measurements map[string]float64
}

func NewRequest(method string, uri string, duration time.Duration, responseCode string) *Request {
	return &Request{method: method, uri: uri, duration: duration, responseCode: responseCode, properties: map[string]string{}, measurements: map[string]float64{}}
}

func (r *Request) AddProperty(key, value string) {
	r.properties[key] = value
}

func (r *Request) AddMeasurement(key string, value float64) {
	r.measurements[key] = value
}

func (r *Request) GetTelemetry() *appinsights.RequestTelemetry {
	request := appinsights.NewRequestTelemetry(r.method, r.uri, r.duration, r.responseCode)
	request.Properties = r.properties
	request.Measurements = r.measurements

	return request
}
