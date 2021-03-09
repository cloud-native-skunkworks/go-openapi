package middleware

import (
	"net/http"
)

func JaegerMiddleWare(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		//spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		//serverSpan := opentracing.GlobalTracer().StartSpan("go-openapi-hop", ext.RPCServerOption(spanCtx))
		//defer serverSpan.Finish()
		handler.ServeHTTP(w, r)
	})
}
