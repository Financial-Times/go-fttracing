package fttracing

import (
	"fmt"
	"os"
	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

const (
	debug = false
	// same span can be set to true for RPC style spans (Zipkin V1) vs Node style (OpenTracing)
	sameSpan = true
	// make Tracer generate 128 bit traceID's for root spans.
	traceID128Bit = true

	zipkinProdEndpoint  = "https://zipkin-upp.in.ft.com/api/v1/spans"
	zipkinStageEndpoint = "https://zipkin-upp-test.in.ft.com/api/v1/spans"

	endpointKey = "endpoint"
)

func constantsForEnv(env string) map[string]string {

	switch env {
	case "prod":
		return map[string]string{
			"endpoint": zipkinProdEndpoint,
		}
	default:
		return map[string]string{
			"endpoint": zipkinStageEndpoint,
		}
	}
}

func NewTracer(serviceName string, hostAndPort string, env string) opentracing.Tracer {

	consts := constantsForEnv(env)

	collector, err := zipkin.NewHTTPCollector(consts[endpointKey])
	if err != nil {
		fmt.Printf("unable to create Zipkin HTTP collector: %+v", err)
		os.Exit(-1)
	}

	// Create our recorder.
	recorder := zipkin.NewRecorder(collector, debug, hostAndPort, serviceName)

	// Create our tracer.
	tracer, err := zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(sameSpan),
		zipkin.TraceID128Bit(traceID128Bit),
	)
	if err != nil {
		fmt.Printf("unable to create Zipkin tracer: %+v", err)
		os.Exit(-1)
	}
	return tracer
}
