package styles

import (
	"github.com/AccentDesign/gcss/props"
)

var (
	white               = props.ColorRGBA(255, 255, 255, 255)
	primary             = props.ColorRGBA(17, 24, 39, 255)
	primaryForeground   = white
	secondary           = props.ColorRGBA(243, 244, 246, 255)
	secondaryForeground = props.ColorRGBA(17, 24, 39, 255)
	destructive         = props.ColorRGBA(220, 38, 38, 255)
	borderColor         = props.ColorRGBA(229, 231, 235, 255)
	mutedForeground     = props.ColorRGBA(75, 85, 99, 255)
	fontSm              = props.UnitRem(0.875)
	leadingTight        = props.UnitRem(1.25)
	radiusSm            = props.UnitRem(0.25)
	radius              = props.UnitRem(0.375)
	size1               = props.UnitRem(0.25)
	size2               = props.UnitRem(0.5)
	size3               = props.UnitRem(0.75)
	size4               = props.UnitRem(1)
	size10              = props.UnitRem(2.5)

	iconChevronDown = props.BackgroundImageURL("data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIxZW0iIGhlaWdodD0iMWVtIiB2aWV3Qm94PSIwIDAgMjQgMjQiPjxwYXRoIGZpbGw9Im5vbmUiIHN0cm9rZT0iY3VycmVudENvbG9yIiBzdHJva2UtbGluZWNhcD0icm91bmQiIHN0cm9rZS1saW5lam9pbj0icm91bmQiIHN0cm9rZS13aWR0aD0iMiIgZD0ibTE5IDlsLTcgN2wtNy03Ii8+PC9zdmc+")
)
