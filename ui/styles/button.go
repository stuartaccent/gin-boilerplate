package styles

import (
	"github.com/AccentDesign/gcss"
	"github.com/AccentDesign/gcss/props"
)

var buttons = []gcss.Style{
	{
		Selector: ".ui-button",
		Props: gcss.Props{
			AlignItems:     props.AlignItemsCenter,
			BorderRadius:   radius,
			Display:        props.DisplayInlineFlex,
			FontSize:       fontSm,
			FontWeight:     props.FontWeightMedium,
			Height:         size10,
			JustifyContent: props.JustifyContentCenter,
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
			BackgroundColor: props.ColorTransparent(),
		},
	},
	{
		Selector: ".ui-button-ghost:hover",
		Props: gcss.Props{
			BackgroundColor: secondary.Alpha(230),
		},
	},
}
