package main

import (
	"errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jessevdk/go-flags"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
	"go-openapi/models"
	logadaptor "go-openapi/pkg/log"
	"go-openapi/restapi"
	"go-openapi/restapi/operations"
	"go-openapi/restapi/operations/user"
	"os"
	"strings"
)

func main() {
	// -----------------------------------------------------------------------------------------------------------------
	// Fake database
	//db, err := storage.LoadLocalDB()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//  -----------------------------------------------------------------------------------------------------------------
	// Logging
	cLogger := log.New()
	cLogger.SetFormatter(&log.JSONFormatter{})
	cLogger.SetOutput(os.Stdout)
	cLogger.SetLevel(log.InfoLevel)
	cLoggerEntry := cLogger.WithFields(log.Fields{
		"app": "go-openapi",
	})
	//  -----------------------------------------------------------------------------------------------------------------
	// Jaeger

	// Configure Jaeger
	cfg := jaegercfg.Configuration{
		ServiceName: "go-openapi",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}
	jMetricsFactory := metrics.NullFactory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(logadaptor.LogrusAdapter{Logger: cLogger}),
		jaegercfg.Metrics(jMetricsFactory))
	if err != nil {
		cLoggerEntry.Error(err.Error())
	}
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		cLoggerEntry.Fatalln(err)
	}

	api := operations.NewGoOpenapiAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	parser := flags.NewParser(server, flags.Default)
	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			cLoggerEntry.Fatalln(err)
		}
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}

	// Example token ---------------------------------------------------------------------------------------------------
	exampleToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
	exampleEmail := "foo@bar.com"
	examplePass := "bar"



	// Stub API code examples ------------------------------------------------------------------------------------------
	api.BearerAuth = func(bearerHeader string) (interface{}, error) {
		bearerToken := strings.Split(bearerHeader, " ")[1]

		if bearerToken == exampleToken {
			return true, nil
		}

		return nil, errors.New("invalid token")
	}
	api.UserLoginHandler = user.LoginHandlerFunc(func(params user.LoginParams) middleware.Responder {

		if *params.Login.Email != exampleEmail && (*params.Login.Password) != examplePass {
			return user.NewLoginNotFound()
		}
		return user.NewLoginOK().WithPayload(&models.LoginSuccess{
			Success: true,
			Token: exampleToken,
		})
	})

	api.UserGetCartHandler = user.GetCartHandlerFunc(func(params user.GetCartParams, i interface{}) middleware.Responder {


		return user.NewGetCartOK()
	})


	// -----------------------------------------------------------------------------------------------------------------
	server.ConfigureAPI()
	server.Port = 8080
	server.Host = "0.0.0.0"
	if err := server.Serve(); err != nil {
		cLoggerEntry.Fatalln(err)
	}
	// -----------------------------------------------------------------------------------------------------------------
}
