package styles

import (
	"github.com/AccentDesign/gcss"
	"github.com/AccentDesign/gcss/props"
)

var menu = []gcss.Style{
	{
		Selector: ".ui-menu",
		Props: gcss.Props{
			BackgroundColor: white,
			Border: props.Border{
				Width: props.UnitPx(1),
				Style: props.BorderStyleSolid,
				Color: borderColor,
			},
			BorderRadius: radius,
			MinWidth:     props.UnitRem(10),
			Overflow:     props.OverflowHidden,
			Padding:      size1,
		},
	},
	{
		Selector: ".ui-menu-label",
		Props: gcss.Props{
			FontSize:      fontSm,
			PaddingTop:    size1,
			PaddingRight:  size2,
			PaddingBottom: size1,
			PaddingLeft:   size2,
		},
	},
	{
		Selector: ".ui-menu-separator",
		Props: gcss.Props{
			BackgroundColor: borderColor,
			Height:          props.UnitPx(1),
			MarginTop:       size1,
			MarginBottom:    size1,
			MarginLeft:      props.UnitRem(-0.25),
			MarginRight:     props.UnitRem(-0.25),
		},
	},
	{
		Selector: ".ui-menu-item",
		Props: gcss.Props{
			BorderRadius:  radiusSm,
			Display:       props.DisplayBlock,
			FontSize:      fontSm,
			FontWeight:    props.FontWeightMedium,
			LineHeight:    leadingTight,
			PaddingTop:    size1,
			PaddingRight:  size2,
			PaddingBottom: size1,
			PaddingLeft:   size2,
		},
	},
	{
		Selector: ".ui-menu-item:hover",
		Props: gcss.Props{
			BackgroundColor: secondary.Alpha(230),
		},
	},
}
