package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	lotel "github.com/rhiadc/grpc_api/client/otel"
	"github.com/rhiadc/grpc_api/client/proto"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

func main() {

	service := os.Getenv("GRPC_CONN_HOST")
	if service == "" {
		service = "grpc_server"
	}
	host := fmt.Sprintf("%s:4040", service)
	conn, err := grpc.Dial(host, grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	client := proto.NewAddServiceClient(conn)

	ctx := context.Background()
	exp, err := lotel.NewExporter(ctx)

	if err != nil {
		log.Fatalf("Error: failed to initialize exporter: %v", err)
	}

	tp := lotel.NewTraceProvider(exp)
	defer func() { _ = tp.Shutdown(ctx) }()

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	tracer := tp.Tracer("restapi", trace.WithInstrumentationVersion("1.0.0"))

	g := gin.Default()
	g.Use(otelgin.Middleware("restapi"))

	g.GET("/add/:a/:b", func(ctx *gin.Context) {
		a, err := strconv.ParseUint(ctx.Param("a"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameter A"})
			return
		}

		b, err := strconv.ParseUint(ctx.Param("b"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameter B"})
			return
		}

		req := &proto.Request{A: int64(a), B: int64(b)}

		if response, err := client.Add(ctx, req); err == nil {
			_, span := tracer.Start(ctx, "add")
			defer span.End()
			span.SetAttributes(
				attribute.String("value", response.String()),
			)
			// setting span as successful
			span.SetStatus(codes.Ok, "Success")
			ctx.JSON(http.StatusOK, gin.H{"result": response.Result})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

	})

	g.GET("/mult/:a/:b", func(ctx *gin.Context) {
		a, err := strconv.ParseUint(ctx.Param("a"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameter A"})
			return
		}

		b, err := strconv.ParseUint(ctx.Param("b"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameter B"})
			return
		}

		req := &proto.Request{A: int64(a), B: int64(b)}

		if response, err := client.Multiply(ctx, req); err == nil {
			tmplName := "user"
			tmplStr := "value {{ .value }}\n"
			tmpl := template.Must(template.New(tmplName).Parse(tmplStr))
			g.SetHTMLTemplate(tmpl)
			otelgin.HTML(ctx, http.StatusOK, tmplName, gin.H{
				"value": response.Result,
			})
			ctx.JSON(http.StatusOK, gin.H{"result": response.Result})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

	})

	if err := g.Run(":8089"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
