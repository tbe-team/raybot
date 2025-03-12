package raybot

import "embed"

//go:embed api/openapi/gen/openapi.yml
var OpenapiSpec []byte

//go:embed ui/dist
var UIFS embed.FS
