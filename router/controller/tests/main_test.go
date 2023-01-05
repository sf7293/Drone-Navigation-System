package tests

import (
	"context"
	"os"
	"testing"

	"github.com/sf7293/Drone-Navigation-System/service/metrix"

	"github.com/opentracing/opentracing-go"

	"github.com/sf7293/Drone-Navigation-System/config"
)

var TestTracingSpan opentracing.Span

func TestMain(m *testing.M) {
	config.Init()

	jeagerConfig := metrix.JaegerConfigsData{
		ServerName:                  config.App.ServerName,
		JaegerAgentAddress:          config.Tracer.JaegerAgentAddress,
		JaegerReporterFlushInterval: config.Tracer.JaegerReporterFlushInterval,
		Testing:                     true,
	}
	metrix.Init(jeagerConfig)
	TestTracingSpan = metrix.CreateSpan("Router.TestTracingSpan", context.Background())
	defer TestTracingSpan.Finish()

	defer func() {
		// close tracer
		_ = metrix.TracerCloser.Close()
	}()

	statusCode := m.Run()
	os.Exit(statusCode)
}
