package main

import (
	"github.com/AlexsJones/go-openapi/restapi/operations/health"
	"os"

	"github.com/AlexsJones/go-openapi/models"
	logadaptor "github.com/AlexsJones/go-openapi/pkg/log"
	"github.com/AlexsJones/go-openapi/pkg/storage"
	"github.com/AlexsJones/go-openapi/restapi"
	"github.com/AlexsJones/go-openapi/restapi/operations"
	"github.com/AlexsJones/go-openapi/restapi/operations/user"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jessevdk/go-flags"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

func main() {
	// -----------------------------------------------------------------------------------------------------------------
	// Fake database
	db, err := storage.LoadLocalDB()
	if err != nil {
		log.Fatal(err)
	}
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

	// Stub API code examples ------------------------------------------------------------------------------------------
	api.HealthGetHealthzHandler = health.GetHealthzHandlerFunc(func(params health.GetHealthzParams) middleware.Responder {
		// Default health check
		return health.NewGetHealthzDefault(200)
	})

	api.UserCreateUserHandler = user.CreateUserHandlerFunc(func(params user.CreateUserParams) middleware.Responder {
		cLoggerEntry.WithFields(log.Fields{
			"username":    params.Body.Username,
			"headers":     params.HTTPRequest.Header,
			"method":      params.HTTPRequest.Method,
			"host":        params.HTTPRequest.Host,
			"requestPath": params.HTTPRequest.RequestURI,
		}).Info("CreateUserHandlerFunc")
		t := opentracing.GlobalTracer()
		span := t.StartSpan("UserCreateUserHandler")
		defer span.Finish()
		tnxGetSpan := tracer.StartSpan(
			"DatabaseGet",
			opentracing.ChildOf(span.Context()),
		)
		defer tnxGetSpan.Finish()
		getTxn := db.Txn(false)
		defer getTxn.Abort()
		raw, err := getTxn.First("user", "username", params.Body.Username)
		if err != nil || raw != nil {
			cLoggerEntry.WithFields(log.Fields{
				"username":    params.Body.Username,
				"headers":     params.HTTPRequest.Header,
				"method":      params.HTTPRequest.Method,
				"host":        params.HTTPRequest.Host,
				"requestPath": params.HTTPRequest.RequestURI,
			}).Warn(err)
			return user.NewCreateUserConflict()
		}
		tnxCreateSpan := tracer.StartSpan(
			"DatabaseCreate",
			opentracing.ChildOf(tnxGetSpan.Context()),
		)
		defer tnxCreateSpan.Finish()
		txn := db.Txn(true)
		defer txn.Abort()
		if err := txn.Insert("user", params.Body); err != nil {
			cLoggerEntry.WithFields(log.Fields{
				"username":    params.Body.Username,
				"headers":     params.HTTPRequest.Header,
				"method":      params.HTTPRequest.Method,
				"host":        params.HTTPRequest.Host,
				"requestPath": params.HTTPRequest.RequestURI,
			}).Warn(err)
		}
		txn.Commit()
		return user.NewCreateUserDefault(201)
	})

	api.UserDeleteUserHandler = user.DeleteUserHandlerFunc(func(params user.DeleteUserParams) middleware.Responder {
		cLoggerEntry.WithFields(log.Fields{
			"username":    params.Username,
			"headers":     params.HTTPRequest.Header,
			"method":      params.HTTPRequest.Method,
			"host":        params.HTTPRequest.Host,
			"requestPath": params.HTTPRequest.RequestURI,
		}).Infof("UserDeleteUserHandler")
		t := opentracing.GlobalTracer()
		span := t.StartSpan("UserDeleteUserHandler")
		defer span.Finish()
		tnxGetSpan := tracer.StartSpan(
			"DatabaseGet",
			opentracing.ChildOf(span.Context()),
		)
		defer tnxGetSpan.Finish()
		getTxn := db.Txn(false)
		defer getTxn.Abort()
		raw, err := getTxn.First("user", "username", params.Username)
		if err != nil || raw == nil {
			cLoggerEntry.WithFields(log.Fields{
				"username":    params.Username,
				"headers":     params.HTTPRequest.Header,
				"method":      params.HTTPRequest.Method,
				"host":        params.HTTPRequest.Host,
				"requestPath": params.HTTPRequest.RequestURI,
			}).Warn(err)
			return user.NewDeleteUserNotFound()
		}
		tnxDelSpan := tracer.StartSpan(
			"DatabaseDelete",
			opentracing.ChildOf(tnxGetSpan.Context()),
		)
		defer tnxDelSpan.Finish()
		writeTxn := db.Txn(true)
		defer writeTxn.Abort()
		if err := writeTxn.Delete("user", raw); err != nil {
			cLoggerEntry.WithFields(log.Fields{
				"username":    params.Username,
				"headers":     params.HTTPRequest.Header,
				"method":      params.HTTPRequest.Method,
				"host":        params.HTTPRequest.Host,
				"requestPath": params.HTTPRequest.RequestURI,
			}).Warn(err)
			return user.NewDeleteUserNotFound()
		}
		writeTxn.Commit()

		return user.NewDeleteUserOK()
	})

	api.UserGetUserByNameHandler = user.GetUserByNameHandlerFunc(func(params user.GetUserByNameParams) middleware.Responder {
		cLoggerEntry.WithFields(log.Fields{
			"username":    params.Username,
			"headers":     params.HTTPRequest.Header,
			"method":      params.HTTPRequest.Method,
			"host":        params.HTTPRequest.Host,
			"requestPath": params.HTTPRequest.RequestURI,
		}).Infof("UserGetUserByNameHandler")
		t := opentracing.GlobalTracer()
		span := t.StartSpan("UserGetUserByNameHandler")
		defer span.Finish()
		tnxSpan := tracer.StartSpan(
			"DatabaseGet",
			opentracing.ChildOf(span.Context()),
		)
		defer tnxSpan.Finish()
		txn := db.Txn(false)
		defer txn.Abort()
		raw, err := txn.First("user", "username", params.Username)
		if err != nil {
			cLoggerEntry.WithFields(log.Fields{
				"username":    params.Username,
				"headers":     params.HTTPRequest.Header,
				"method":      params.HTTPRequest.Method,
				"host":        params.HTTPRequest.Host,
				"requestPath": params.HTTPRequest.RequestURI,
			}).Warn(err)
		}
		if raw == nil {
			return user.NewGetUserByNameNotFound()
		}
		u := raw.(*models.User)
		return user.NewGetUserByNameOK().WithPayload(u)
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
