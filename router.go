package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Router struct {
	mux        *chi.Mux
	spec       *openapi3.Swagger
	specRouter *openapi3filter.Router
}

func (r *Router) Use(middlewares ...func(http.Handler) http.Handler) {
	r.mux.Use(middlewares...)
}

func (r *Router) ListenAndServe(addr string) {
	http.ListenAndServe(addr, r.mux)
}

func (router *Router) mockHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route, pathParams, err := router.specRouter.FindRoute(r.Method, r.URL)

		ctx := context.TODO()
		if err != nil {
			log.Fatalf("Error finding route: %s", err)
			os.Exit(1)
			// panic(err)
		}
		// Validate request
		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    r,
			PathParams: pathParams,
			Route:      route,
		}
		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			panic(err)
		}

		var (
			respStatus      = 200
			respContentType = "application/json"
			respBody        = bytes.NewBufferString(`{}`)
		)

		log.Println("Response:", respStatus)
		responseValidationInput := &openapi3filter.ResponseValidationInput{
			RequestValidationInput: requestValidationInput,
			Status:                 respStatus,
			Header: http.Header{
				"Content-Type": []string{
					respContentType,
				},
			},
		}
		if respBody != nil {
			data, _ := json.Marshal(respBody)
			responseValidationInput.SetBodyBytes(data)
		}

		// Validate response.
		if err := openapi3filter.ValidateResponse(ctx, responseValidationInput); err != nil {
			panic(err)
		}
	}
}

func (r *Router) initHandler() {
	methods := [9]string{
		"GET",
		"HEAD",
		"POST",
		"PUT",
		"PATCH", // RFC 5789
		"DELETE",
		"CONNECT",
		"OPTIONS",
		"TRACE",
	}

	for _, method := range methods {
		r.mux.Method(method, "/*", r.mockHandler())
	}

}

func NewRouter() *Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	return &Router{
		mux: r,
	}
}

func (r *Router) WithSpec(spec *openapi3.Swagger) *Router {
	r.spec = spec
	r.specRouter = openapi3filter.NewRouter().WithSwagger(r.spec)
	r.initHandler()
	return r
}
