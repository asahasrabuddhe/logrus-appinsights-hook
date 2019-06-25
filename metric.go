package logrusApplicationInsightsHook

import "github.com/Microsoft/ApplicationInsights-Go/appinsights"

type Metric struct {
	name       string
	value      float64
	sessionId  string
	properties map[string]string
}

func (m *Metric) AddProperty(key, value string) {
	m.properties[key] = value
}

func (m *Metric) SetSessionId(sessionId string) {
	m.sessionId = sessionId
}

func (m *Metric) GetTelemetry() *appinsights.MetricTelemetry {
	metric := appinsights.NewMetricTelemetry(m.name, m.value)
	metric.Properties = m.properties

	if m.sessionId != "" {
		metric.Tags.Session().SetId(m.sessionId)
	}

	return metric
}
