package html

import (
	"context"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin/render"
	"net/http"
)

type Render struct {
	Fallback render.HTMLRender
}

func (r *Render) Instance(s string, d any) render.Render {
	if component, ok := d.(templ.Component); ok {
		return &TemplRender{
			Ctx:       context.Background(),
			Status:    -1,
			Component: component,
		}
	}
	return r.Fallback.Instance(s, d)
}

type TemplRender struct {
	Ctx       context.Context
	Status    int
	Component templ.Component
}

func (r TemplRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	if r.Status != -1 {
		w.WriteHeader(r.Status)
	}
	if r.Component != nil {
		return r.Component.Render(r.Ctx, w)
	}
	return nil
}

func (r TemplRender) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}
