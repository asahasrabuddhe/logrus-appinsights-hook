package logrusApplicationInsightsHook

import "github.com/Microsoft/ApplicationInsights-Go/appinsights"

type Metric struct {
	name       string
	value      float64
	properties map[string]string
}

func (m *Metric) AddProperty(key, value string) {
	m.properties[key] = value
}

func (m *Metric) GetTelemetry(sessionId string) *appinsights.MetricTelemetry {
	metric := appinsights.NewMetricTelemetry(m.name, m.value)
	metric.Properties = m.properties

	if sessionId != "" {
		metric.Tags.Session().SetId(sessionId)
	}

	return metric
}
