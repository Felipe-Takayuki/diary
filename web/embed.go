package web

import _ "embed"

// IndexHTML contains the embedded single-page application frontend.
//
//go:embed template/index.html
var IndexHTML []byte

// FaviconSVG contains the embedded vector favicon.
//
//go:embed static/favicon.svg
var FaviconSVG []byte

// FaviconICO contains the embedded multi-resolution ICO favicon.
//
//go:embed static/favicon.ico
var FaviconICO []byte
