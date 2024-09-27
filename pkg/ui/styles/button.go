package styles

import (
	"github.com/AccentDesign/gcss"
	"github.com/AccentDesign/gcss/props"
	"github.com/AccentDesign/gcss/variables"
)

func (ss *StyleSheet) Buttons() Styles {
	return Styles{
		{
			Selector: ".ui-button",
			Props: gcss.Props{
				AlignItems:     props.AlignItemsCenter,
				BorderRadius:   variables.Size1H,
				Display:        props.DisplayInlineFlex,
				FontSize:       variables.Size3H,
				FontWeight:     props.FontWeightMedium,
				Height:         variables.Size10,
				JustifyContent: props.JustifyContentCenter,
				LineHeight:     variables.Size5,
				PaddingTop:     variables.Size2,
				PaddingRight:   variables.Size4,
				PaddingBottom:  variables.Size2,
				PaddingLeft:    variables.Size4,
			},
		},
		{
			Selector: ".ui-button:hover",
			Props: gcss.Props{
				Cursor: props.CursorPointer,
			},
		},
	}
}

func (t *Theme) Buttons() Styles {
	return Styles{
		{
			Selector: ".ui-button-primary",
			Props: gcss.Props{
				BackgroundColor: t.Primary,
				Color:           t.PrimaryForeground,
			},
		},
		{
			Selector: ".ui-button-primary:hover",
			Props: gcss.Props{
				BackgroundColor: t.Primary.Alpha(230),
			},
		},
		{
			Selector: ".ui-button-secondary",
			Props: gcss.Props{
				BackgroundColor: t.Secondary,
				Color:           t.SecondaryForeground,
			},
		},
		{
			Selector: ".ui-button-secondary:hover",
			Props: gcss.Props{
				BackgroundColor: t.Secondary.Alpha(230),
			},
		},
		{
			Selector: ".ui-button-ghost",
			Props: gcss.Props{
				BackgroundColor: props.ColorTransparent(),
			},
		},
		{
			Selector: ".ui-button-ghost:hover",
			Props: gcss.Props{
				BackgroundColor: t.Secondary.Alpha(230),
			},
		},
	}
}
