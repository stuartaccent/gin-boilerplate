package styles

import (
	"github.com/AccentDesign/gcss"
	"io"
	"sort"
	"sync"
)

type StyleCache struct {
	styles map[string]gcss.Style
	sync.Mutex
}

var allStyles = [][]gcss.Style{buttons, form, menu}

// New creates a new StyleCache with the default styles.
func New() *StyleCache {
	cache := &StyleCache{
		styles: make(map[string]gcss.Style),
	}
	for _, styles := range allStyles {
		for _, style := range styles {
			cache.AddStyle(style)
		}
	}
	return cache
}

// SortedKeys returns the keys of the styles in sorted order.
func (sc *StyleCache) SortedKeys() []string {
	keys := make([]string, len(sc.styles))
	i := 0
	for key := range sc.styles {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

// AddStyle adds a style to the cache.
func (sc *StyleCache) AddStyle(style gcss.Style) {
	sc.Lock()
	defer sc.Unlock()
	sc.styles[style.Selector] = style
}

// WriteCss writes the CSS to the writer.
func (sc *StyleCache) WriteCss(w io.Writer) error {
	sc.Lock()
	defer sc.Unlock()
	for _, key := range sc.SortedKeys() {
		style := sc.styles[key]
		if err := style.CSS(w); err != nil {
			return err
		}
	}
	return nil
}
