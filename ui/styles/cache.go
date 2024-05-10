package styles

import (
	"bytes"
	"github.com/AccentDesign/gcss"
	"io"
	"sort"
	"sync"
)

type StyleCache struct {
	styles map[string]gcss.Style
	css    map[string][]byte
	sync.Mutex
}

var allStyles = [][]gcss.Style{buttons, form, menu}

// New creates a new StyleCache with the default styles.
func New() (*StyleCache, error) {
	cache := &StyleCache{
		styles: make(map[string]gcss.Style),
		css:    make(map[string][]byte),
	}
	for _, styles := range allStyles {
		for _, style := range styles {
			if err := cache.AddStyle(style); err != nil {
				return nil, err
			}
		}
	}
	return cache, nil
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
func (sc *StyleCache) AddStyle(style gcss.Style) error {
	sc.Lock()
	defer sc.Unlock()
	sc.styles[style.Selector] = style
	var b bytes.Buffer
	if err := style.CSS(&b); err != nil {
		return err
	}
	sc.css[style.Selector] = b.Bytes()
	return nil
}

// WriteCss writes the CSS to the writer.
func (sc *StyleCache) WriteCss(w io.Writer) error {
	sc.Lock()
	defer sc.Unlock()
	for _, key := range sc.SortedKeys() {
		css := sc.css[key]
		if _, err := w.Write(css); err != nil {
			return err
		}
	}
	return nil
}
