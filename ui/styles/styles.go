package styles

type Size string
type Variant string

const (
	Md      Size    = "md"
	Primary Variant = "primary"
)

var (
	Error       = []string{"text-red-500", "text-sm"}
	Input       = []string{"bg-gray-50", "border", "border-gray-300", "text-gray-900", "text-sm", "rounded-lg", "focus:ring-blue-500", "focus:border-blue-500", "block", "w-full", "p-2.5"}
	Label       = []string{"text-sm", "font-medium", "mb-1"}
	Menu        = []string{"z-10", "font-normal", "bg-white", "divide-y", "divide-gray-100", "rounded-lg", "shadow", "absolute"}
	MenuItems   = []string{"py-2", "text-sm", "text-gray-700"}
	MenuLink    = []string{"block", "px-4", "py-2", "hover:bg-gray-200"}
	btnVariants = map[Variant][]string{
		Primary: {"text-white", "bg-slate-900", "hover:bg-slate-950", "focus:ring-blue-500"},
	}
	btnSizes = map[Size][]string{
		Md: {"font-medium", "rounded-lg", "text-sm", "px-5", "py-2.5"},
	}
)

func Button(v Variant, s Size) []string {
	return append(append([]string{}, btnVariants[v]...), btnSizes[s]...)
}
