package styles

import (
	"github.com/AccentDesign/gcss"
	"github.com/AccentDesign/gcss/props"
	"github.com/AccentDesign/gcss/variables"
)

func (ss *StyleSheet) Menu() Styles {
	return Styles{
		{
			Selector: ".ui-menu",
			Props: gcss.Props{
				Border: props.Border{
					Width: props.UnitPx(1),
					Style: props.BorderStyleSolid,
				},
				BorderRadius: variables.Size1H,
				MinWidth:     props.UnitRem(10),
				Overflow:     props.OverflowHidden,
				Padding:      variables.Size1,
			},
		},
		{
			Selector: ".ui-menu-label",
			Props: gcss.Props{
				FontSize:      variables.Size3H,
				PaddingTop:    variables.Size1,
				PaddingRight:  variables.Size2,
				PaddingBottom: variables.Size1,
				PaddingLeft:   variables.Size2,
			},
		},
		{
			Selector: ".ui-menu-separator",
			Props: gcss.Props{
				Height:       props.UnitPx(1),
				MarginTop:    variables.Size1,
				MarginBottom: variables.Size1,
				MarginLeft:   props.UnitRem(-0.25),
				MarginRight:  props.UnitRem(-0.25),
			},
		},
		{
			Selector: ".ui-menu-item",
			Props: gcss.Props{
				BorderRadius:  variables.Size1,
				Display:       props.DisplayBlock,
				FontSize:      variables.Size3H,
				FontWeight:    props.FontWeightMedium,
				LineHeight:    variables.Size5,
				PaddingTop:    variables.Size1,
				PaddingRight:  variables.Size2,
				PaddingBottom: variables.Size1,
				PaddingLeft:   variables.Size2,
			},
		},
	}
}

func (t *Theme) Menu() Styles {
	return Styles{
		{
			Selector: ".ui-menu",
			Props: gcss.Props{
				BackgroundColor: t.Background,
				BorderColor:     t.Border,
			},
		},
		{
			Selector: ".ui-menu-separator",
			Props: gcss.Props{
				BackgroundColor: t.Border,
			},
		},
		{
			Selector: ".ui-menu-item:hover",
			Props: gcss.Props{
				BackgroundColor: t.Secondary.Alpha(230),
			},
		},
	}
}
