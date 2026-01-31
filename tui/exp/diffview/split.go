package diffview

import (
	"github.com/aymanbagabas/go-udiff"
)

// SplitHunk represents a hunk in split (side-by-side) view format.
type SplitHunk struct {
	FromLine int
	ToLine   int
	Lines    []*SplitLine
}

// SplitLine contains before/after line pairs for split view.
type SplitLine struct {
	Before *udiff.Line
	After  *udiff.Line
}

// HunkToSplit converts a unified diff hunk to split format.
func HunkToSplit(h *udiff.Hunk) SplitHunk {
	lines := make([]udiff.Line, len(h.Lines))
	copy(lines, h.Lines)

	sh := SplitHunk{
		FromLine: h.FromLine,
		ToLine:   h.ToLine,
		Lines:    make([]*SplitLine, 0, len(lines)),
	}

	for len(lines) > 0 {
		ul := lines[0]
		lines = lines[1:]

		var sl SplitLine

		switch ul.Kind {
		// For equal lines, add as is
		case udiff.Equal:
			sl.Before = &ul
			sl.After = &ul

		// For inserted lines, set after and keep before as nil
		case udiff.Insert:
			sl.Before = nil
			sl.After = &ul

		// For deleted lines, set before and loop over the next lines
		// searching for the equivalent after line.
		case udiff.Delete:
			sl.Before = &ul

			// Search for matching insert line
			for i, l := range lines {
				if l.Kind == udiff.Insert {
					sl.After = &l
					// Remove the insert line from lines
					lines = append(lines[:i], lines[i+1:]...)
					break
				}
				if l.Kind == udiff.Equal {
					break
				}
			}
		}

		sh.Lines = append(sh.Lines, &sl)
	}

	return sh
}
