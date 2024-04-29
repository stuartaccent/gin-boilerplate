package styles

import (
	"github.com/AccentDesign/gcss"
	"github.com/AccentDesign/gcss/props/border"
	"github.com/AccentDesign/gcss/props/display"
	"github.com/AccentDesign/gcss/props/font"
	"github.com/AccentDesign/gcss/props/overflow"
	"github.com/AccentDesign/gcss/props/unit"
)

var menu = []gcss.Style{
	{
		Selector: ".ui-menu",
		Props: gcss.Props{
			BackgroundColor: white,
			Border: border.Border{
				Width: unit.Px(1),
				Style: border.StyleSolid,
				Color: borderColor,
			},
			BorderRadius: radius,
			MinWidth:     unit.Rem(10),
			Overflow:     overflow.Hidden,
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
			Height:          unit.Px(1),
			MarginTop:       size1,
			MarginBottom:    size1,
			MarginLeft:      unit.Rem(-0.25),
			MarginRight:     unit.Rem(-0.25),
		},
	},
	{
		Selector: ".ui-menu-item",
		Props: gcss.Props{
			BorderRadius:  radiusSm,
			Display:       display.Block,
			FontSize:      fontSm,
			FontWeight:    font.WeightMedium,
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
