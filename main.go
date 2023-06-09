package main

import (
	"os"

	l "go-jeager/logger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func main() {
	router := gin.Default()

	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(os.Getenv("JEAGER_URL"))))
	if err != nil {
		log.Error().Err(err)
	}

	// Create a resource
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(os.Getenv("JEAGER_SERVICE")),
	)

	// Create a tracer provider with the Jaeger exporter
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// Register the tracer provider as the global provider
	otel.SetTracerProvider(tp)

	// Create a Gin router
	// Middleware to start a new span for each request
	router.Use(l.JaegerMiddleware)

	router.GET("/hello", func(c *gin.Context) {
		// span := trace.SpanFromContext(c.Request.Context())
		// defer span.End()

		//Start --> Example code for Child Span

		// // if span.SpanContext().IsValid() {
		// // Create a child span
		// _, childSpan := otel.Tracer(c.Request.URL.Path).Start(
		// 	c.Request.Context(),
		// 	c.Request.URL.Path,
		// 	trace.WithLinks(trace.Link{SpanContext: span.SpanContext()}),
		// 	trace.WithSpanKind(trace.SpanKindClient),
		// )
		// defer childSpan.End()

		//end --> Example code for Child Span

		l.JSON(c, 400, gin.H{
			"message": "Hello, World!",
		})
	})

	port := os.Getenv("SERVICE_PORT")
	log.Info().Msg("Starting server on :" + port)
	router.Run(":" + port)
}
