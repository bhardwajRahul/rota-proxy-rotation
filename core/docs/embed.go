package docs

import _ "embed"

// SwaggerJSON holds the OpenAPI spec embedded into the binary at build time.
//
// It is served by the API's /api/v1/swagger.json route. Embedding (instead of
// reading docs/swagger.json from disk at runtime) means the docs render in any
// deployment regardless of the process working directory — notably the Docker
// image, whose final stage ships only the compiled binary and no docs/ folder,
// which previously made the documentation page render empty (issue #20).
//
//go:embed swagger.json
var SwaggerJSON []byte
