package logrusApplicationInsightsHook

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
	"github.com/Microsoft/ApplicationInsights-Go/appinsights/contracts"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	TelemetryType       string = "app_insights_telemetry_type"
	EventTelemetry      string = "app_insights_event_telemetry"
	MetricTelemetry     string = "app_insights_metric_telemetry"
	RequestTelemetry    string = "app_insights_request_telemetry"
	DependencyTelemetry string = "app_insights_dependency_telemetry"
)

const (
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
	levels         []logrus.Level
	excludedFields map[string]struct{}
	filters        map[string]func(interface{}) interface{}
}

func (a *ApplicationInsightsHook) Fire(e *logrus.Entry) error {
	if telemetryType, ok := e.Data[TelemetryType]; ok {
		switch telemetryType {
		case EventTelemetry:
			if et, ok := e.Data[EventField].(*Event); ok {
				if sessionId, ok := e.Data["session_id"].(string); ok {
					et.SetSessionId(sessionId)
				}

				a.client.Track(et.GetTelemetry())
			} else {
				return errors.New("event telemetry expects an event field")
			}
		case MetricTelemetry:
			if mt, ok := e.Data[MetricField].(*Metric); ok {
				if sessionId, ok := e.Data["session_id"].(string); ok {
					mt.SetSessionId(sessionId)
				}

				a.client.Track(mt.GetTelemetry())
			} else {
				return errors.New("metric telemetry expects an metric field")
			}
		case RequestTelemetry:
			if rt, ok := e.Data[RequestField].(*Request); ok {
				if sessionId, ok := e.Data["session_id"].(string); ok {
					rt.SetSessionId(sessionId)
				}

				a.client.Track(rt.GetTelemetry())
			} else {
				return errors.New("request telemetry expects an request field")
			}
		case DependencyTelemetry:
			if dt, ok := e.Data[DependencyField].(*Dependency); ok {
				if sessionId, ok := e.Data["session_id"].(string); ok {
					dt.SetSessionId(sessionId)
				}

				a.client.Track(dt.GetTelemetry())
			} else {
				return errors.New("dependency telemetry expects an dependency field")
			}

		default:
			return errors.New("invalid telemetry type defined")
		}

		delete(e.Data, telemetryType)

	} else {
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
			trace := appinsights.NewTraceTelemetry(e.Message, levels[e.Level])
			for key, val := range e.Data {
				if sessionId, ok := e.Data["session_id"].(string); ok {
					trace.Tags.Session().SetId(sessionId)
					e.Data["session_id"] = sessionId
					delete(e.Data, key)
				} else {
					fVal, _ := formatData(val)
					trace.Properties[key] = fVal
				}
			}
			a.client.Track(trace)
		}
	}

	return nil
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
	Role               string
	Version            string
	Debug              bool
}

func NewApplicationInsightsHook(opts ApplicationInsightsHookOpts) (*ApplicationInsightsHook, error) {
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

	if opts.Role != "" {
		client.Context().Tags.Cloud().SetRole(opts.Role)
	}

	if opts.Version != "" {
		client.Context().Tags.Application().SetVer(opts.Version)
	}

	if opts.Debug {
		appinsights.NewDiagnosticsMessageListener(func(msg string) error {
			fmt.Printf("[%s] %s\n", time.Now().Format(time.UnixDate), msg)
			return nil
		})
	}

	return &ApplicationInsightsHook{
		client:         client,
		levels:         logrusLevels,
		excludedFields: make(map[string]struct{}),
		filters:        make(map[string]func(interface{}) interface{}),
	}, nil
}

// formatData returns value as a suitable format.
func formatData(value interface{}) (string, error) {
	switch value := value.(type) {
	case json.Marshaler:
		if bytes, err := value.MarshalJSON(); err != nil {
			return string(bytes), nil
		} else {
			return "", err
		}
	case error:
		return value.Error(), nil
	case fmt.Stringer:
		return value.String(), nil
	default:
		return fmt.Sprintf("%v", value), nil
	}
}
