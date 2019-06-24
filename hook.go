package logrusApplicationInsightsHook

import (
	"errors"
	"fmt"
	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
	"github.com/Microsoft/ApplicationInsights-Go/appinsights/contracts"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	SessionIdField  string = "app_insights_session_id"
	EventField      string = "app_insights_event"
	MetricField     string = "app_insights_metric"
	RequestField    string = "app_insights_request"
	DependencyField string = "app_insights_dependency"
)

var logrusLevels = []logrus.Level{
	logrus.PanicLevel,
	logrus.FatalLevel,
	logrus.ErrorLevel,
	logrus.WarnLevel,
	logrus.InfoLevel,
	logrus.DebugLevel,
	logrus.TraceLevel,
}

var levels = map[logrus.Level]contracts.SeverityLevel{
	logrus.PanicLevel: contracts.Critical,
	logrus.FatalLevel: contracts.Error,
	logrus.ErrorLevel: contracts.Error,
	logrus.WarnLevel:  contracts.Warning,
	logrus.InfoLevel:  contracts.Information,
	logrus.DebugLevel: contracts.Information,
	logrus.TraceLevel: contracts.Verbose,
}

type ApplicationInsightsHook struct {
	client         appinsights.TelemetryClient
	async          bool
	levels         []logrus.Level
	excludedFields map[string]struct{}
	filters        map[string]func(interface{}) interface{}
}

func (a *ApplicationInsightsHook) Fire(e *logrus.Entry) error {
	if e.Level == logrus.ErrorLevel || e.Level == logrus.FatalLevel || e.Level == logrus.PanicLevel {
		if val, ok := e.Data["error"].(error); ok {
			a.client.TrackException(val)
		} else {
			a.client.TrackException(errors.New(e.Message))
		}

		if e.Level == logrus.FatalLevel || e.Level == logrus.PanicLevel {
			<-a.client.Channel().Close()
		}

	} else {
		for key, val := range e.Data {
			switch key {
			case SessionIdField:
				trace := appinsights.NewTraceTelemetry(e.Message, levels[e.Level])
				if sessionId, ok := val.(string); ok {
					trace.Tags.Session().SetId(sessionId)
					e.Data["session_id"] = sessionId
				}
				a.client.Track(trace)
			case EventField:
				if event, ok := val.(*Event); ok {
					if val, ok := e.Data["session_id"]; ok {

					}
					a.client.Track(event.GetTelemetry())
				}
			case MetricField:
				if metric, ok := val.(*Metric); ok {
					if val, ok := e.Data["session_id"]; ok {

					}
					a.client.Track(metric.GetTelemetry())
				}
				delete(e.Data, MetricField)
			case RequestField:
				if request, ok := val.(*Request); ok {
					if val, ok := e.Data["session_id"]; ok {

					}
					a.client.Track(request.GetTelemetry())
				}
			case DependencyField:
				if dependency, ok := val.(*Dependency); ok {
					if val, ok := e.Data["session_id"]; ok {

					}
					a.client.Track(dependency.GetTelemetry())
				}
			}

			delete(e.Data, key)
		}
	}
	return nil
}

func (a *ApplicationInsightsHook) Async() bool {
	return a.async
}

func (a *ApplicationInsightsHook) SetAsync(async bool) {
	a.async = async
}

func (a *ApplicationInsightsHook) Levels() []logrus.Level {
	return a.levels
}

func (a *ApplicationInsightsHook) SetLevels(levels []logrus.Level) {
	a.levels = levels
}

func (a *ApplicationInsightsHook) Close() {
	<-a.client.Channel().Close()
}

type ApplicationInsightsHookOpts struct {
	InstrumentationKey string
	EndpointUrl        string
	MaxBatchSize       int
	MaxBatchInterval   time.Duration
}

func NewApplicationInsightsHook(name string, opts ApplicationInsightsHookOpts) (*ApplicationInsightsHook, error) {
	if opts.InstrumentationKey == "" {
		return nil, errors.New("instrumentation key is not provided")
	}

	config := appinsights.NewTelemetryConfiguration(opts.InstrumentationKey)

	if opts.EndpointUrl != "" {
		config.EndpointUrl = opts.EndpointUrl
	}

	if opts.MaxBatchSize != 0 {
		config.MaxBatchSize = opts.MaxBatchSize
	}

	if opts.MaxBatchInterval != 0 {
		config.MaxBatchInterval = opts.MaxBatchInterval
	}

	client := appinsights.NewTelemetryClientFromConfig(config)

	if name != "" {
		client.Context().Tags.Cloud().SetRole(name)
	}

	appinsights.NewDiagnosticsMessageListener(func(msg string) error {
		fmt.Printf("[%s] %s\n", time.Now().Format(time.UnixDate), msg)
		return nil
	})

	return &ApplicationInsightsHook{
		client:         client,
		levels:         logrusLevels,
		excludedFields: make(map[string]struct{}),
		filters:        make(map[string]func(interface{}) interface{}),
	}, nil
}
