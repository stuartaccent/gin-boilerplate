package styles

type buttonStyles struct {
	Base     string
	Variants S
	Sizes    S
}

var buttons = buttonStyles{
	Base: "inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:opacity-50 disabled:pointer-events-none ring-offset-background",
	Variants: S{
		"primary":     "bg-primary text-primary-foreground hover:bg-primary/90",
		"secondary":   "bg-secondary text-secondary-foreground hover:bg-secondary/80",
		"destructive": "bg-destructive text-destructive-foreground hover:bg-destructive/90",
		"outline":     "border border-input hover:bg-accent hover:text-accent-foreground",
		"ghost":       "hover:bg-accent hover:text-accent-foreground",
		"link":        "underline-offset-4 hover:underline text-primary",
	},
	Sizes: S{
		"sm": "h-9 px-3",
		"md": "h-10 px-4",
		"lg": "h-11 px-8",
	},
}

func Button(variant, size string) string {
	classes := buttons.Base

	if cls, ok := buttons.Variants[variant]; ok {
		classes += " " + cls
	} else {
		classes += " " + buttons.Variants["primary"]
	}
	if cls, ok := buttons.Sizes[size]; ok {
		classes += " " + cls
	} else {
		classes += " " + buttons.Sizes["md"]
	}

	return classes
}
