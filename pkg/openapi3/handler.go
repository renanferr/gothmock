package openapi3

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Handler struct {
	mux        *chi.Mux
	spec       *openapi3.Swagger
	specRouter *openapi3filter.Router
	options    *HandlerOptions
}

func (h *Handler) Use(middlewares ...func(http.Handler) http.Handler) {
	h.mux.Use(middlewares...)
}

func (h *Handler) Get() http.Handler {
	return h.mux
}

func (h *Handler) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(h.spec.Paths)
		route, pathParams, err := h.specRouter.FindRoute(r.Method, r.URL)

		ctx := context.TODO()
		if err != nil {
			err = fmt.Errorf("route not found: %s %s: %w", r.Method, r.URL, err)
			log.Println(err)
			notFoundResponse(w)
			return
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

		op := route.PathItem.GetOperation(r.Method)
		log.Println(op)

		resp := op.Responses.Get(h.options.ResponseStatus)
		mediaType := resp.Value.Content.Get(h.options.ResponseContent)
		respBody, err := json.Marshal(mediaType.Example)
		if err != nil {
			panic(err)
		}

		w.Write(respBody)
		w.Header().Set("Content-type", h.options.ResponseContent)
		w.WriteHeader(h.options.ResponseStatus)

		return
	}
}

func (r *Handler) initializeHandlers() {
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
		r.mux.Method(method, "/*", r.handler())
	}

}

func NewHandler() *Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	g := &Handler{
		mux:     r,
		options: &HandlerOptions{},
	}
	g.initializeHandlers()
	return g
}

func (h *Handler) WithOptions(options *HandlerOptions) *Handler {
	h.options = options
	return h
}

func (r *Handler) WithSpec(spec *openapi3.Swagger) *Handler {
	r.spec = spec
	r.specRouter = openapi3filter.NewRouter().WithSwagger(r.spec)
	return r
}

func notFoundResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(http.StatusText(http.StatusNotFound)))
}
