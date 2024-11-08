package middleware

import (
	"net/http"
)

// HTMX contains utilities for dealing with HTMX requests
type HTMX struct {
	Request  *http.Request
	Response http.ResponseWriter
}

// IsHTMXRequest checks if the request is an HTMX request.
func (h *HTMX) IsHTMXRequest() bool {
	return h.Request.Header.Get("HX-Request") == "true"
}

// SetRedirect sets the `HX-Redirect` header for HTMX requests.
func (h *HTMX) SetRedirect(url string) {
	h.Response.Header().Set("HX-Redirect", url)
}

// SetPushUrl sets the `HX-Push-Url` header for HTMX requests.
func (h *HTMX) SetPushUrl(url string) {
	h.Response.Header().Set("HX-Push-Url", url)
}

// SetRefresh sets the `HX-Refresh` header for HTMX requests.
func (h *HTMX) SetRefresh() {
	h.Response.Header().Set("HX-Refresh", "true")
}

// SetTrigger sets the `HX-Trigger` header for HTMX requests.
func (h *HTMX) SetTrigger(content string) {
	h.Response.Header().Set("HX-Trigger", content)
}

// SetTriggerAfterSettle sets the `HX-Trigger-After-Settle` header for HTMX requests.
func (h *HTMX) SetTriggerAfterSettle(content string) {
	h.Response.Header().Set("HX-Trigger-After-Settle", content)
}

// SetTriggerAfterSwap sets the `HX-Trigger-After-Swap` header for HTMX requests.
func (h *HTMX) SetTriggerAfterSwap(content string) {
	h.Response.Header().Set("HX-Trigger-After-Swap", content)
}
