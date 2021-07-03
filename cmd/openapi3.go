package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	swagger "github.com/getkin/kin-openapi/openapi3"
	"github.com/renanferr/gothmock/pkg/openapi3"
	"github.com/spf13/cobra"
)

var (
	openApi3Cmd = &cobra.Command{
		Use:   "openapi3 [OpenAPI3 file path]",
		Short: "Mocks an API",
		Long:  `Mocks an API from an OpenAPI 3 specification file`,
		Args:  cobra.MinimumNArgs(1),
		Run:   MockOpenApi3,
	}
)

func init() {
	rootCmd.AddCommand(openApi3Cmd)
}

func MockOpenApi3(cmd *cobra.Command, args []string) {
	spec, err := LoadOpenApi3Spec(args[0])
	if err != nil {
		log.Fatal(fmt.Errorf("error loading OpenAPI3 specification file: %w", err))
	}

	handler := openapi3.NewHandler().
		WithSpec(spec).
		WithOptions(openapi3.NewHandlerOptions(status, content))

	ListenAndServe(handler.Get(), GetAddr())

}

func ListenAndServe(handler http.Handler, addr string) {
	go func() {
		log.Printf("starting server at http://0.0.0.0%s", addr)
		log.Fatal(http.ListenAndServe(addr, handler))
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
}

func GetAddr() string {
	if !strings.HasPrefix(addr, ":") {
		addr = ":" + addr
	}
	return addr
}

func LoadOpenApi3Spec(filePath string) (*swagger.Swagger, error) {
	if strings.HasPrefix(filePath, "http") {
		return loadOpenApi3SpecFromURI(filePath)
	}
	return swagger.NewSwaggerLoader().LoadSwaggerFromFile(filePath)
}

func loadOpenApi3SpecFromURI(uri string) (*swagger.Swagger, error) {
	loc, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	return swagger.NewSwaggerLoader().LoadSwaggerFromURI(loc)
}
