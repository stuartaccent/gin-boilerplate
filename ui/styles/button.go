package styles

import (
	"github.com/AccentDesign/gcss"
	"github.com/AccentDesign/gcss/props/align"
	"github.com/AccentDesign/gcss/props/colors"
	"github.com/AccentDesign/gcss/props/display"
	"github.com/AccentDesign/gcss/props/font"
	"github.com/AccentDesign/gcss/props/justify"
)

var buttons = []gcss.Style{
	{
		Selector: ".ui-button",
		Props: gcss.Props{
			AlignItems:     align.ItemsCenter,
			BorderRadius:   radius,
			Display:        display.InlineFlex,
			FontSize:       fontSm,
			FontWeight:     font.WeightMedium,
			Height:         size10,
			JustifyContent: justify.ContentCenter,
			LineHeight:     leadingTight,
			PaddingTop:     size2,
			PaddingRight:   size4,
			PaddingBottom:  size2,
			PaddingLeft:    size4,
		},
	},
	{
		Selector: ".ui-button-primary",
		Props: gcss.Props{
			BackgroundColor: primary,
			Color:           primaryForeground,
		},
	},
	{
		Selector: ".ui-button-primary:hover",
		Props: gcss.Props{
			BackgroundColor: primary.Alpha(230),
		},
	},
	{
		Selector: ".ui-button-secondary",
		Props: gcss.Props{
			BackgroundColor: secondary,
			Color:           secondaryForeground,
		},
	},
	{
		Selector: ".ui-button-secondary:hover",
		Props: gcss.Props{
			BackgroundColor: secondary.Alpha(230),
		},
	},
	{
		Selector: ".ui-button-ghost",
		Props: gcss.Props{
			BackgroundColor: colors.Transparent(),
		},
	},
	{
		Selector: ".ui-button-ghost:hover",
		Props: gcss.Props{
			BackgroundColor: secondary.Alpha(230),
		},
	},
}
