package middleware

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"net/http"
)

func JaegerMiddleWare(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		serverSpan := opentracing.GlobalTracer().StartSpan(r.RequestURI, ext.RPCServerOption(spanCtx))
		defer serverSpan.Finish()

		ext.SpanKindRPCClient.Set(serverSpan)
		ext.HTTPUrl.Set(serverSpan, r.URL.String())
		ext.HTTPMethod.Set(serverSpan, r.Method)
		// Inject the span context into the headers
		opentracing.GlobalTracer().Inject(serverSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		handler.ServeHTTP(w, r)
	})
}
