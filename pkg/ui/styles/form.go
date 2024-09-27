package styles

import (
	"github.com/AccentDesign/gcss"
	"github.com/AccentDesign/gcss/props"
	"github.com/AccentDesign/gcss/variables"
)

func (ss *StyleSheet) Form() Styles {
	return Styles{
		{
			Selector: ".ui-input,.ui-select",
			Props: gcss.Props{
				Border: props.Border{
					Width: props.UnitPx(1),
					Style: props.BorderStyleSolid,
				},
				BorderRadius:  variables.Size1H,
				Display:       props.DisplayFlex,
				FontSize:      variables.Size3H,
				Height:        variables.Size10,
				LineHeight:    variables.Size5,
				PaddingTop:    variables.Size2,
				PaddingRight:  variables.Size3,
				PaddingBottom: variables.Size2,
				PaddingLeft:   variables.Size3,
				Width:         props.UnitPercent(100),
			},
		},
		{
			Selector: ".ui-input::file-selector-button",
			Props: gcss.Props{
				BackgroundColor: props.ColorTransparent(),
				BorderWidth:     props.UnitRaw(0),
				FontSize:        variables.Size3H,
				FontWeight:      props.FontWeightMedium,
			},
		},
		{
			Selector: ".ui-input-label",
			Props: gcss.Props{
				FontSize:   variables.Size3H,
				FontWeight: props.FontWeightMedium,
				LineHeight: variables.Size5,
			},
		},
		{
			Selector: ".ui-input-help",
			Props: gcss.Props{
				FontSize: variables.Size3H,
			},
		},
		{
			Selector: ".ui-input-error",
			Props: gcss.Props{
				FontSize: variables.Size3H,
			},
		},
		{
			Selector: ".ui-select:where(:not([size]))",
			Props: gcss.Props{
				Appearance:       props.AppearanceNone,
				PaddingRight:     variables.Size10,
				PrintColorAdjust: props.PrintColorAdjustExact,
				BackgroundImage:  iconChevronDown,
				BackgroundPosition: props.BackgroundPositionEdges(
					props.BackgroundPositionEdge{Position: props.BackgroundPositionRight, Unit: variables.Size3},
					props.BackgroundPositionEdge{Position: props.BackgroundPositionCenter},
				),
				BackgroundRepeat: props.BackgroundRepeatNoRepeat,
				BackgroundSize:   props.BackgroundSizeDimension(props.UnitEm(1), props.UnitEm(1)),
			},
		},
	}
}

func (t *Theme) Form() Styles {
	return Styles{
		{
			Selector: ".ui-input",
			Props: gcss.Props{
				BackgroundColor: t.Background,
				BorderColor:     t.Border,
			},
		},
		{
			Selector: ".ui-input-help",
			Props: gcss.Props{
				Color: t.MutedForeground,
			},
		},
		{
			Selector: ".ui-input-error",
			Props: gcss.Props{
				Color: t.Destructive,
			},
		},
		{
			Selector: ".ui-select",
			Props: gcss.Props{
				BackgroundColor: t.Background,
				BorderColor:     t.Border,
			},
		},
	}
}
