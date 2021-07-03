package cmd

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	swagger "github.com/getkin/kin-openapi/openapi3"
	"github.com/renanferr/gothmock/pkg/openapi3"
	"github.com/spf13/cobra"
)

const (
	DefaultContent = "application/json"
	DefaultStatus  = http.StatusOK
	DefaultAddr    = ":8085"
)

var (
	content string
	status  int
	addr    string

	openApi3Cmd = &cobra.Command{
		Use:   "openapi3 [OpenAPI3 file path]",
		Short: "Mocks an API",
		Long:  `Mocks an API from an OpenAPI 3 specification file`,
		Args:  cobra.MinimumNArgs(1),
		Run:   MockOpenApi3,
	}
)

func init() {
	openApi3Cmd.Flags().StringVar(&content, "content", DefaultContent, "Response Content-type")
	openApi3Cmd.Flags().StringVar(&addr, "port", DefaultAddr, "Server listening port")
	openApi3Cmd.Flags().IntVar(&status, "status", DefaultStatus, "Response Status")
	rootCmd.AddCommand(openApi3Cmd)
}

func MockOpenApi3(cmd *cobra.Command, args []string) {
	filepath := args[0]
	swagger, err := swagger.NewSwaggerLoader().LoadSwaggerFromFile(filepath)
	if err != nil {
		log.Fatalf("Error loading swagger file: %s", err)
	}

	handler := openapi3.NewHandler().
		WithSpec(swagger).
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
