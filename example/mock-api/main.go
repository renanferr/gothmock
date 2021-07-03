package main

import (
	"log"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/renanferr/gothmock/pkg/http"
)

func main() {
	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("./swagger.yml")
	if err != nil {
		log.Fatalf("Error loading swagger file: %s", err)
	}
	for _, s := range swagger.Servers {
		log.Println(s.URL)
	}

	r := http.NewRouter().WithSpec(swagger)

	r.ListenAndServe(":8085")
}
