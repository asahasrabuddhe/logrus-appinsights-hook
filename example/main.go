package main

import (
	logrusApplicationInsightsHook "github.com/asahasrabuddhe/logrus-appinsights"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

func main() {
	defer recoverFunc()

	opts := logrusApplicationInsightsHook.ApplicationInsightsHookOpts{
		InstrumentationKey: "instrumentation key",
		MaxBatchSize:       10,
		MaxBatchInterval:   2 * time.Second,
	}

	hook, err := logrusApplicationInsightsHook.NewApplicationInsightsHook("test 2", opts)
	if err != nil {
		log.Fatal(err)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logrus.SetLevel(logrus.TraceLevel)

	logrus.AddHook(hook)

	lg := logrus.WithField(logrusApplicationInsightsHook.SessionIdField, "s1234")

	lg.Trace("trace")
	lg.Debug("debug")

	sampleEvent := logrusApplicationInsightsHook.NewEvent("sample event")
	sampleEvent.AddProperty("property", "value")

	lg.WithField(logrusApplicationInsightsHook.EventField, sampleEvent).Info("info")

	lg.Warn("warn")
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
