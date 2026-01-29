package layout

import "image"

// FlexChild represents a child in a flex layout with its constraint.
type FlexChild struct {
	// Constraint defines the size constraint for this child
	Constraint Constraint
	// Grow indicates if this child should grow to fill remaining space
	Grow bool
	// Shrink indicates if this child can shrink below its constraint size
	Shrink bool
}

// NewFlexChild creates a new FlexChild with the given constraint.
func NewFlexChild(constraint Constraint) FlexChild {
	return FlexChild{
		Constraint: constraint,
		Grow:       false,
		Shrink:     false,
	}
}

// WithGrow enables grow behavior for the FlexChild.
func (c FlexChild) WithGrow() FlexChild {
	c.Grow = true
	return c
}

// WithShrink enables shrink behavior for the FlexChild.
func (c FlexChild) WithShrink() FlexChild {
	c.Shrink = true
	return c
}

// RowLayout divides the area horizontally into child areas.
// Children are arranged from left to right.
//
// The layout algorithm:
// 1. Calculate fixed sizes from non-growing children
// 2. Calculate remaining space
// 3. Distribute remaining space to growing children equally
// 4. Assign areas from left to right
//
// Example:
//
//	children := []layout.FlexChild{
//	    layout.NewFlexChild(layout.Fixed(20)),       // 20 columns fixed
//	    layout.NewFlexChild(layout.Grow{}).WithGrow(), // Fills remaining
//	    layout.NewFlexChild(layout.Fixed(10)),        // 10 columns fixed
//	}
//	areas := RowLayout(area, children)
func RowLayout(area Area, children []FlexChild) []Area {
	return flexLayout(area, children, true)
}

// ColumnLayout divides the area vertically into child areas.
// Children are arranged from top to bottom.
//
// The layout algorithm is the same as RowLayout, but vertical.
//
// Example:
//
//	children := []layout.FlexChild{
//	    layout.NewFlexChild(layout.Fixed(5)),         // 5 rows fixed
//	    layout.NewFlexChild(layout.Grow{}).WithGrow(), // Fills remaining
//	    layout.NewFlexChild(layout.Percent(20)),      // 20% of remaining
//	}
//	areas := ColumnLayout(area, children)
func ColumnLayout(area Area, children []FlexChild) []Area {
	return flexLayout(area, children, false)
}

// flexLayout is the common implementation for both row and column layouts.
func flexLayout(area Area, children []FlexChild, horizontal bool) []Area {
	if len(children) == 0 {
		return []Area{}
	}

	if len(children) == 1 && children[0].Grow {
		return []Area{area}
	}

	available := getAvailableSize(area, horizontal)

	// Calculate sizes for all children
	sizes := calculateChildSizes(available, children, horizontal)

	// Create child areas
	areas := make([]Area, len(children))
	currentPos := getStartPos(area, horizontal)

	for i, size := range sizes {
		if horizontal {
			areas[i] = calculateHorizontalArea(area, currentPos.X, size)
			currentPos.X += size
		} else {
			areas[i] = calculateVerticalArea(area, currentPos.Y, size)
			currentPos.Y += size
		}
	}

	return areas
}

// getAvailableSize returns the available width or height based on horizontal flag.
func getAvailableSize(area Area, horizontal bool) int {
	if horizontal {
		return area.Dx()
	}
	return area.Dy()
}

// getStartPos returns the starting position based on horizontal flag.
func getStartPos(area Area, horizontal bool) image.Point {
	if horizontal {
		return area.Rect().Min
	}
	return area.Rect().Min
}

// calculateChildSizes calculates the size for each child based on constraints.
func calculateChildSizes(available int, children []FlexChild, horizontal bool) []int {
	sizes := make([]int, len(children))

	// First pass: calculate fixed sizes
	var fixedTotal int
	growCount := 0

	for i, child := range children {
		if !child.Grow {
			sizes[i] = child.Constraint.Apply(available)
			fixedTotal += sizes[i]
		} else {
			sizes[i] = 0
			growCount++
		}
	}

	// Second pass: distribute remaining space to growing children
	if growCount > 0 {
		remaining := available - fixedTotal
		if remaining > 0 {
			growSize := remaining / growCount
			for i, child := range children {
				if child.Grow {
					sizes[i] = growSize
				}
			}
		}
	}

	return sizes
}

// calculateHorizontalArea creates a horizontal child area.
func calculateHorizontalArea(parent Area, startX, width int) Area {
	return Area{
		Min: image.Pt(startX, parent.Rect().Min.Y),
		Max: image.Pt(startX+width, parent.Rect().Max.Y),
	}
}

// calculateVerticalArea creates a vertical child area.
func calculateVerticalArea(parent Area, startY, height int) Area {
	return Area{
		Min: image.Pt(parent.Rect().Min.X, startY),
		Max: image.Pt(parent.Rect().Max.X, startY+height),
	}
}

// FlexRow is a convenience function to create a row layout with inline children.
// See RowLayout for details.
//
// Example:
//
//	areas := layout.FlexRow(area,
//	    layout.Fixed(20),
//	    layout.Grow{},
//	    layout.Percent(30),
//	)
func FlexRow(area Area, constraints ...Constraint) []Area {
	children := make([]FlexChild, len(constraints))
	for i, c := range constraints {
		children[i] = NewFlexChild(c)
	}
	return RowLayout(area, children)
}

// FlexColumn is a convenience function to create a column layout with inline children.
// See ColumnLayout for details.
//
// Example:
//
//	areas := layout.FlexColumn(area,
//	    layout.Fixed(5),
//	    layout.Grow{},
//	    layout.Percent(20),
//	)
func FlexColumn(area Area, constraints ...Constraint) []Area {
	children := make([]FlexChild, len(constraints))
	for i, c := range constraints {
		children[i] = NewFlexChild(c)
	}
	return ColumnLayout(area, children)
}
