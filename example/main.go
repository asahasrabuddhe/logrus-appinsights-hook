package main

import (
	"github.com/sirupsen/logrus"
	logrusApplicationInsightsHook "go.ajitem.com/logrus-appinsights-hook"
	"log"
	"time"
)

func main() {
	defer recoverFunc()

	opts := logrusApplicationInsightsHook.ApplicationInsightsHookOpts{
		InstrumentationKey: "instrumentation key",
		MaxBatchSize:       10,
		MaxBatchInterval:   2 * time.Second,
		Role:               "myapptest",
		Version:            "0.1.2-staging",
		Debug:              true,
	}

	hook, err := logrusApplicationInsightsHook.NewApplicationInsightsHook(opts)
	if err != nil {
		log.Fatal(err)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logrus.SetLevel(logrus.TraceLevel)

	logrus.AddHook(hook)

	lg := logrus.WithField("session_id", "s1234")

	lg.Trace("trace")
	lg.Debug("debug")

	sampleEvent := logrusApplicationInsightsHook.NewEvent("sample event")
	sampleEvent.AddProperty("property", "value")

	lg.WithFields(logrus.Fields{
		logrusApplicationInsightsHook.TelemetryType: logrusApplicationInsightsHook.EventTelemetry,
		logrusApplicationInsightsHook.EventField:    sampleEvent,
	}).Info("")

	lgg := logrus.WithFields(logrus.Fields{
		"field1": 1,
		"field2": "two",
		"field3": 3.3,
	})

	lgg.Warn("warn")
	//lg.Error("error")
	//lg.Fatal("fatal")
	//lg.Panic("panic")

	hook.Close()
}

func recoverFunc() {
	if e := recover(); e != nil {
		log.Println("recovered from", e)
	}
}
