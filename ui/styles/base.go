package styles

import (
	"github.com/AccentDesign/gcss"
	"github.com/AccentDesign/gcss/variables"
)

func (ss *StyleSheet) Base() Styles {
	return Styles{
		{
			Selector: "body",
			Props: gcss.Props{
				MinHeight: variables.FullScreenHeight,
			},
			CustomProps: []gcss.CustomProp{
				{"-webkit-font-smoothing", "antialiased"},
				{"-moz-osx-font-smoothing", "grayscale"},
			},
		},
	}
}

func (t *Theme) Base() Styles {
	return Styles{
		{
			Selector: "body",
			Props: gcss.Props{
				BackgroundColor: t.Background,
				Color:           t.Foreground,
			},
		},
	}
}
