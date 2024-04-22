package styles

type inputStyles struct {
	Base  string
	Sizes S
}

var (
	InputError = "text-red-500 text-sm"
	InputLabel = "text-sm font-medium leading-none"
)

var inputs = inputStyles{
	Base: "flex w-full rounded-md border border-input bg-background text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",
	Sizes: S{
		"md": "h-10 py-2 px-4",
	},
}

func Input(size string) string {
	classes := inputs.Base

	if cls, ok := inputs.Sizes[size]; ok {
		classes += " " + cls
	} else {
		classes += " " + inputs.Sizes["md"]
	}

	return classes
}
