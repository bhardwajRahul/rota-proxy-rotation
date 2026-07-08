package handlers

import (
	"fmt"
	"net/http"

	scalar "github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/alpkeskin/rota/core/docs"
)

// DocumentationHandler handles API documentation endpoints
type DocumentationHandler struct{}

// NewDocumentationHandler creates a new DocumentationHandler
func NewDocumentationHandler() *DocumentationHandler {
	return &DocumentationHandler{}
}

// ServeDocumentation serves the API documentation using Scalar.
// The spec is passed as embedded content (not a URL) so Scalar never fetches it
// server-side — that fetch used the external request host and failed behind the
// reverse proxy, where the core can't reach its own public address (#20).
func (h *DocumentationHandler) ServeDocumentation(w http.ResponseWriter, r *http.Request) {
	htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecContent: string(docs.SwaggerJSON),
		CustomOptions: scalar.CustomOptions{
			PageTitle: "Rota Proxy API Documentation",
		},
		DarkMode: true,
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating documentation: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintln(w, htmlContent)
}
