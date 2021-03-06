package logrusApplicationInsightsHook

import "github.com/Microsoft/ApplicationInsights-Go/appinsights"

type Event struct {
	name         string
	sessionId    string
	properties   map[string]string
	measurements map[string]float64
}

func NewEvent(name string) *Event {
	return &Event{name: name, properties: make(map[string]string)}
}

func (e *Event) AddProperty(key, value string) {
	e.properties[key] = value
}

func (e *Event) AddMeasurement(key string, value float64) {
	e.measurements[key] = value
}

func (e *Event) SetSessionId(sessionId string) {
	e.sessionId = sessionId
}

func (e *Event) GetTelemetry() *appinsights.EventTelemetry {
	event := appinsights.NewEventTelemetry(e.name)
	event.Properties = e.properties
	event.Measurements = e.measurements

	if e.sessionId != "" {
		event.Tags.Session().SetId(e.sessionId)
	}

	return event
}
