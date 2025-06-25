// Copyright 2019 DeepMap, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package codegen

import (
	"bufio"
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/oapi-codegen/oapi-codegen/v2/pkg/util"
)

// ...existing code up to GetResponseTypeDefinitions...
func (o *OperationDefinition) GetResponseTypeDefinitions() ([]ResponseTypeDefinition, error) {
	var tds []ResponseTypeDefinition

	if o.Spec == nil || o.Spec.Responses == nil {
		return tds, nil
	}

	sortedResponsesKeys := SortedMapKeys(o.Spec.Responses.Map())
	for _, responseName := range sortedResponsesKeys {
		responseRef := o.Spec.Responses.Value(responseName)

		// We can only generate a type if we have a value:
		if responseRef.Value != nil {
			jsonCount := 0
			for mediaType := range responseRef.Value.Content {
				if util.IsMediaTypeJson(mediaType) {
					jsonCount++
				}
			}

			sortedContentKeys := SortedMapKeys(responseRef.Value.Content)
			for _, contentTypeName := range sortedContentKeys {
				contentType := responseRef.Value.Content[contentTypeName]
				// We can only generate a type if we have a schema:
				if contentType.Schema != nil {
					responseSchema, err := GenerateGoSchema(contentType.Schema, []string{o.OperationId, responseName})
					if err != nil {
						return nil, fmt.Errorf("unable to determine Go type for %s.%s: %w", o.OperationId, contentTypeName, err)
					}

					var typeName string
					switch {

					// HAL+JSON:
					case StringInArray(contentTypeName, contentTypesHalJSON):
						typeName = fmt.Sprintf("HALJSON%s", nameNormalizer(responseName))
					case contentTypeName == "application/json":
						// if it's the standard application/json
						typeName = fmt.Sprintf("JSON%s", nameNormalizer(responseName))
					// Vendored JSON
					case StringInArray(contentTypeName, contentTypesJSON) || util.IsMediaTypeJson(contentTypeName):
						baseTypeName := fmt.Sprintf("%s%s", nameNormalizer(contentTypeName), nameNormalizer(responseName))

						typeName = strings.ReplaceAll(baseTypeName, "Json", "JSON")
					// YAML:
					case StringInArray(contentTypeName, contentTypesYAML):
						typeName = fmt.Sprintf("YAML%s", nameNormalizer(responseName))
					// XML:
					case StringInArray(contentTypeName, contentTypesXML):
						typeName = fmt.Sprintf("XML%s", nameNormalizer(responseName))
					default:
						continue
					}

					td := ResponseTypeDefinition{
						TypeDefinition: TypeDefinition{
							TypeName: typeName,
							Schema:   responseSchema,
						},
						ResponseName:              responseName,
						ContentTypeName:           contentTypeName,
						AdditionalTypeDefinitions: responseSchema.GetAdditionalTypeDefs(),
					}
					if IsGoTypeReference(responseRef.Ref) {
						refType, err := RefPathToGoType(responseRef.Ref)
						if err != nil {
							return nil, fmt.Errorf("error dereferencing response Ref: %w", err)
						}
						if jsonCount > 1 && util.IsMediaTypeJson(contentTypeName) {
							refType += mediaTypeToCamelCase(contentTypeName)
						}
						// Patch: append JSONResponse suffix for JSON content types to match v2.0.0 behavior
						if util.IsMediaTypeJson(contentTypeName) {
							refType += "JSONResponse"
						}
						td.Schema.RefType = refType
					}
					tds = append(tds, td)
				}
			}
		}
	}
	return tds, nil
}
// ...rest of the file unchanged...
