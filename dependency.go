package logrusApplicationInsightsHook

import (
	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
	"time"
)

type Dependency struct {
	name           string
	dependencyType string
	target         string
	duration       time.Duration
	success        bool
	sessionId      string
	properties     map[string]string
	measurements   map[string]float64
}

func NewDependency(name string, dependencyType string, target string, duration time.Duration, success bool) *Dependency {
	return &Dependency{name: name, dependencyType: dependencyType, target: target, duration: duration, success: success}
}

func (d *Dependency) AddProperty(key, value string) {
	d.properties[key] = value
}

func (d *Dependency) AddMeasurement(key string, value float64) {
	d.measurements[key] = value
}

func (d *Dependency) SetSessionId(sessionId string) {
	d.sessionId = sessionId
}

func (d *Dependency) GetTelemetry() *appinsights.RemoteDependencyTelemetry {
	dependency := appinsights.NewRemoteDependencyTelemetry(d.name, d.dependencyType, d.target, d.success)
	dependency.Duration = d.duration

	if d.sessionId != "" {
		dependency.Tags.Session().SetId(d.sessionId)
	}

	return dependency
}
