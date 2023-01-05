package metrix

import (
	"context"
	"io"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/sf7293/Drone-Navigation-System/utils/logger"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

var (
	Tracer       opentracing.Tracer
	TracerCloser io.Closer
)

type JaegerConfigsData struct {
	ServerName                  string
	Testing                     bool
	JaegerReporterFlushInterval int
	JaegerAgentAddress          string
}

func Init(inputCfg JaegerConfigsData) {
	defer func() {
		_ = logger.ZSLogger.Sync()
	}()
	var err error

	jaegerReporterConfig, disableTracer, SamplerConfig := createTracerConfig(inputCfg)
	cfg := jaegercfg.Configuration{
		ServiceName: inputCfg.ServerName,
		Disabled:    disableTracer,
		Sampler:     SamplerConfig,
		Reporter:    jaegerReporterConfig,
	}

	TracerLogger := jaegerlog.StdLogger
	TracerMetricsFactory := metrics.NullFactory

	Tracer, TracerCloser, err = cfg.NewTracer(
		jaegercfg.Logger(TracerLogger),
		jaegercfg.Metrics(TracerMetricsFactory),
	)
	if err != nil {
		logger.ZSLogger.Errorf("error on creating tracer: %s", err)
	}
}

func createTracerConfig(JaegerConfigsData JaegerConfigsData) (*jaegercfg.ReporterConfig, bool, *jaegercfg.SamplerConfig) {
	var jaegerReporterConfig *jaegercfg.ReporterConfig
	var jaegerSamplerConfig *jaegercfg.SamplerConfig
	var disableTracer bool
	if JaegerConfigsData.Testing {
		disableTracer = true
		jaegerReporterConfig = &jaegercfg.ReporterConfig{
			LogSpans: false,
		}
		jaegerSamplerConfig = &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		}
		return jaegerReporterConfig, disableTracer, jaegerSamplerConfig
	}
	jaegerReporterConfig = &jaegercfg.ReporterConfig{
		LogSpans:            true,
		BufferFlushInterval: time.Duration(JaegerConfigsData.JaegerReporterFlushInterval) * time.Second,
		LocalAgentHostPort:  JaegerConfigsData.JaegerAgentAddress,
	}
	jaegerSamplerConfig = &jaegercfg.SamplerConfig{
		Type:  jaeger.SamplerTypeConst,
		Param: 1,
	}
	return jaegerReporterConfig, disableTracer, jaegerSamplerConfig
}

func CreateSpan(SpanName string, c context.Context) opentracing.Span {
	var parentCtx opentracing.SpanContext
	ctx := opentracing.SpanFromContext(c)
	if ctx != nil {
		parentCtx = ctx.Context()
	}
	return Tracer.StartSpan(SpanName, opentracing.ChildOf(parentCtx))

}

func CreateChildSpan(SpanName string, ParentContext opentracing.Span) opentracing.Span {
	return Tracer.StartSpan(SpanName, opentracing.ChildOf(ParentContext.Context()))
}

func GetSpanFromGin(spanName string, c *gin.Context) opentracing.Span {
	parentSpan := opentracing.SpanFromContext(c.Request.Context())
	return StartSpan(spanName, opentracing.ChildOf(parentSpan.Context()))
}

func StartSpan(name string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return Tracer.StartSpan(name, opts...)
}
