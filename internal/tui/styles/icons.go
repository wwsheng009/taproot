package styles

const (
	CheckIcon         string = "✓"
	ErrorIcon         string = "×"
	WarningIcon       string = "⚠"
	InfoIcon          string = "ⓘ"
	HintIcon          string = "∵"
	SpinnerIcon       string = "..."
	ArrowRightIcon    string = "→"
	CenterSpinnerIcon string = "⋯"
	LoadingIcon       string = "⟳"
	ImageIcon         string = "■"
	TextIcon          string = "☰"
	ModelIcon         string = "◇"

	// Tool call icons
	ToolPending string = "●"
	ToolSuccess string = "✓"
	ToolError   string = "×"

	BorderThin  string = "│"
	BorderThick string = "▌"

	// Todo icons
	TodoCompletedIcon string = "✓"
	TodoPendingIcon   string = "•"
)

var SelectionIgnoreIcons = []string{
	BorderThin,
	BorderThick,
}
