package codegen

import (
	"go/format"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This test ensures that referenced response components for JSON content types
// are generated with the 'JSONResponse' suffix for backward compatibility.
func TestReferencedResponseType_JSONResponseSuffix(t *testing.T) {
	spec := `openapi: 3.0.1
info:
  title: Example API
  version: 1.0.0
paths:
  /foo:
    get:
      operationId: getFoo
      responses:
        400:
          $ref: '#/components/responses/400'
components:
  responses:
    400:
      description: Bad Request
      content:
        application/json:
          schema:
            type: object
            properties:
              message:
                type: string
            required:
              - message
`
	packageName := "testrefresp"
	opts := Configuration{
		PackageName: packageName,
		Generate: GenerateOptions{
			Strict: true,
			Models: true,
		},
	}
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	swagger, err := loader.LoadFromData([]byte(spec))
	require.NoError(t, err)

	code, err := Generate(swagger, opts)
	assert.NoError(t, err)
	assert.NotEmpty(t, code)

	_, err = format.Source([]byte(code))
	assert.NoError(t, err)

	// The referenced response type for 400 should have the JSONResponse suffix
	assert.Contains(t, code, "type GetFoo400JSONResponse struct")
}
