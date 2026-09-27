package web

import _ "embed"

// IndexHTML contains the embedded single-page application frontend.
//
//go:embed template/index.html
var IndexHTML []byte
