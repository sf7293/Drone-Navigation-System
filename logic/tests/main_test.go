package tests

import (
	"context"
	"os"

	"github.com/opentracing/opentracing-go"
	"github.com/sf7293/Drone-Navigation-System/service/metrix"

	"testing"

	//"github.com/opentracing/opentracing-go"
	"github.com/sf7293/Drone-Navigation-System/config"
	"github.com/sf7293/Drone-Navigation-System/constants"
)

var TestTracingSpan opentracing.Span

func TestMain(m *testing.M) {
	config.App.Mode = constants.TestConfigMode
	config.Init()

	jeagerConfig := metrix.JaegerConfigsData{
		ServerName:                  config.App.ServerName,
		JaegerAgentAddress:          config.Tracer.JaegerAgentAddress,
		JaegerReporterFlushInterval: config.Tracer.JaegerReporterFlushInterval,
		Testing:                     true,
	}
	metrix.Init(jeagerConfig)
	TestTracingSpan = metrix.CreateSpan("Logic.TestTracingSpan", context.Background())

	defer TestTracingSpan.Finish()

	defer func() {
		// close tracer
		_ = metrix.TracerCloser.Close()
	}()

	//repository.Init()

	statusCode := m.Run()
	os.Exit(statusCode)
}
