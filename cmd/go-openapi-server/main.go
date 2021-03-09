package main

import (
	"github.com/AlexsJones/go-openapi/models"
	"github.com/AlexsJones/go-openapi/pkg/storage"
	"github.com/AlexsJones/go-openapi/restapi"
	"github.com/AlexsJones/go-openapi/restapi/operations"
	"github.com/AlexsJones/go-openapi/restapi/operations/user"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jessevdk/go-flags"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"log"
	"os"
)

func main() {

	// Fake database
	db, err := storage.LoadLocalDB()
	if err != nil {
		log.Fatal(err)
	}
	// Configure Jaeger
	cfg := jaegercfg.Configuration{
		ServiceName: "go-openapi",
		Sampler:     &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter:    &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		jLogger.Error(err.Error())
	}
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewGoOpenapiAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	parser := flags.NewParser(server, flags.Default)
	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
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

	// Stub API code examples ------------------------------------------------------------------------------------------
	api.UserCreateUserHandler = user.CreateUserHandlerFunc(func(params user.CreateUserParams) middleware.Responder {
		t := opentracing.GlobalTracer()
		span := t.StartSpan("UserCreateUserHandler")
		defer span.Finish()
		tnxSpan := tracer.StartSpan(
			"DatabaseTx",
			opentracing.ChildOf(span.Context()),
		)
		defer tnxSpan.Finish()
		txn := db.Txn(true)
		if err := txn.Insert("user",params.Body); err != nil {
			log.Fatal(err)
		}
		txn.Commit()
		return user.NewCreateUserDefault(201)
	})
	api.UserGetUserByNameHandler = user.GetUserByNameHandlerFunc(func(params user.GetUserByNameParams) middleware.Responder {
		t := opentracing.GlobalTracer()
		span := t.StartSpan("UserGetUserByNameHandler")
		defer span.Finish()
		tnxSpan := tracer.StartSpan(
			"DatabaseTx",
			opentracing.ChildOf(span.Context()),
		)
		defer tnxSpan.Finish()
		txn := db.Txn(true)
		raw, err := txn.First("user", "username",params.Username)
		if err != nil {
			panic(err)
		}
		txn.Commit()
		if raw == nil {
			return user.NewGetUserByNameNotFound()
 		}
		u := raw.(*models.User)
		return user.NewGetUserByNameOK().WithPayload(u)
	})
	// -----------------------------------------------------------------------------------------------------------------
	server.ConfigureAPI()
	server.Port = 8080
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
	// -----------------------------------------------------------------------------------------------------------------

}
