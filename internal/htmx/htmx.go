package htmx

import (
	"net/http"
)

// Helper contains utilities for dealing with HTMX requests
type Helper struct {
	Request  *http.Request
	Response http.ResponseWriter
}

// New creates a new HTMX helper.
func New(req *http.Request, res http.ResponseWriter) *Helper {
	return &Helper{
		Request:  req,
		Response: res,
	}
}

// IsHTMXRequest checks if the request is an HTMX request.
func (h *Helper) IsHTMXRequest() bool {
	return h.Request.Header.Get("HX-Request") == "true"
}

// SetRedirect sets the `HX-Redirect` header for HTMX requests.
func (h *Helper) SetRedirect(url string) {
	h.Response.Header().Set("HX-Redirect", url)
}

// SetPushUrl sets the `HX-Push-Url` header for HTMX requests.
func (h *Helper) SetPushUrl(url string) {
	h.Response.Header().Set("HX-Push-Url", url)
}

// SetRefresh sets the `HX-Refresh` header for HTMX requests.
func (h *Helper) SetRefresh() {
	h.Response.Header().Set("HX-Refresh", "true")
}

// SetTrigger sets the `HX-Trigger` header for HTMX requests.
func (h *Helper) SetTrigger(content string) {
	h.Response.Header().Set("HX-Trigger", content)
}

// SetTriggerAfterSettle sets the `HX-Trigger-After-Settle` header for HTMX requests.
func (h *Helper) SetTriggerAfterSettle(content string) {
	h.Response.Header().Set("HX-Trigger-After-Settle", content)
}

// SetTriggerAfterSwap sets the `HX-Trigger-After-Swap` header for HTMX requests.
func (h *Helper) SetTriggerAfterSwap(content string) {
	h.Response.Header().Set("HX-Trigger-After-Swap", content)
}
