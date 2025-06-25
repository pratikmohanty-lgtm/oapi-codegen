// ...existing code...
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
// ...existing code...