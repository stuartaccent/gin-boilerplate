package styles

import (
	"bytes"
	"github.com/AccentDesign/gcss"
	"github.com/AccentDesign/gcss/variables"
	"io"
	"slices"
	"sync"
)

type (
	Styles     []gcss.Style
	StyleSheet struct {
		Themes []*Theme
		css    bytes.Buffer
		mutex  sync.Mutex
	}
)

func (ss *StyleSheet) CSS(w io.Writer) error {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	if ss.css.Len() > 0 {
		_, err := w.Write(ss.css.Bytes())
		return err
	}

	for _, style := range slices.Concat(
		ss.Resets(),
		ss.Base(),
		ss.Buttons(),
		ss.Form(),
		ss.Menu(),
	) {
		if err := style.CSS(&ss.css); err != nil {
			return err
		}
	}
	for _, theme := range ss.Themes {
		if err := theme.CSS(&ss.css); err != nil {
			return err
		}
	}

	_, err := w.Write(ss.css.Bytes())
	return err
}

func NewStyleSheet() *StyleSheet {
	return &StyleSheet{
		Themes: []*Theme{
			{
				MediaQuery:          "@media (prefers-color-scheme: light)",
				Background:          variables.White,
				Foreground:          variables.Neutral900,
				Primary:             variables.Neutral900,
				PrimaryForeground:   variables.White,
				Secondary:           variables.Neutral300,
				SecondaryForeground: variables.Neutral900,
				Destructive:         variables.Red700,
				Border:              variables.Neutral300,
				MutedForeground:     variables.Neutral500,
			},
			{
				MediaQuery:          "@media (prefers-color-scheme: dark)",
				Background:          variables.Neutral900,
				Foreground:          variables.Neutral100,
				Primary:             variables.White,
				PrimaryForeground:   variables.Neutral900,
				Secondary:           variables.Neutral600,
				SecondaryForeground: variables.Neutral900,
				Destructive:         variables.Red600,
				Border:              variables.Neutral600,
				MutedForeground:     variables.Neutral400,
			},
		},
	}
}
