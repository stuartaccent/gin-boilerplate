package htmx

import (
	"net/http"
)

// HTMXHelper contains utilities for dealing with HTMX requests
type HTMXHelper struct {
	Request  *http.Request
	Response http.ResponseWriter
}

// IsHTMXRequest checks if the request is an HTMX request.
func (h *HTMXHelper) IsHTMXRequest() bool {
	return h.Request.Header.Get("HX-Request") == "true"
}

// SetRedirect sets the `HX-Redirect` header for HTMX requests.
func (h *HTMXHelper) SetRedirect(url string) {
	h.Response.Header().Set("HX-Redirect", url)
}

// SetPushUrl sets the `HX-Push-Url` header for HTMX requests.
func (h *HTMXHelper) SetPushUrl(url string) {
	h.Response.Header().Set("HX-Push-Url", url)
}

// SetRefresh sets the `HX-Refresh` header for HTMX requests.
func (h *HTMXHelper) SetRefresh() {
	h.Response.Header().Set("HX-Refresh", "true")
}

// SetTrigger sets the `HX-Trigger` header for HTMX requests.
func (h *HTMXHelper) SetTrigger(content string) {
	h.Response.Header().Set("HX-Trigger", content)
}

// SetTriggerAfterSettle sets the `HX-Trigger-After-Settle` header for HTMX requests.
func (h *HTMXHelper) SetTriggerAfterSettle(content string) {
	h.Response.Header().Set("HX-Trigger-After-Settle", content)
}

// SetTriggerAfterSwap sets the `HX-Trigger-After-Swap` header for HTMX requests.
func (h *HTMXHelper) SetTriggerAfterSwap(content string) {
	h.Response.Header().Set("HX-Trigger-After-Swap", content)
}
