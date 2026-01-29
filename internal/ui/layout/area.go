package layout

import "image"

// Position represents a coordinate in 2D space.
type Position struct {
	X int
	Y int
}

// Area represents a rectangular area in 2D space.
// Use image.Rectangle methods for common operations.
type Area image.Rectangle

// NewArea creates a new Area with the given bounds.
func NewArea(minX, minY, maxX, maxY int) Area {
	return Area(image.Rect(minX, minY, maxX, maxY))
}

// TopLeft returns the minimum (top-left) position of the area.
func (a Area) TopLeft() Position {
	return Position{X: a.Rect().Min.X, Y: a.Rect().Min.Y}
}

// BottomRight returns the maximum (bottom-right) position of the area.
func (a Area) BottomRight() Position {
	return Position{X: a.Rect().Max.X, Y: a.Rect().Max.Y}
}

// Dx returns the width of the area.
func (a Area) Dx() int {
	return a.Rect().Dx()
}

// Dy returns the height of the area.
func (a Area) Dy() int {
	return a.Rect().Dy()
}

// Empty returns true if the area has zero width or height.
func (a Area) Empty() bool {
	return a.Rect().Empty()
}

// Rect returns the underlying image.Rectangle.
func (a Area) Rect() image.Rectangle {
	return image.Rectangle(a)
}

// Intersect returns the intersection of this area with another.
func (a Area) Intersect(other Area) Area {
	return Area(a.Rect().Intersect(other.Rect()))
}

// Union returns the union of this area with another.
func (a Area) Union(other Area) Area {
	return Area(a.Rect().Union(other.Rect()))
}

// Constraint represents a size constraint for layout calculations.
type Constraint interface {
	// Apply applies the constraint to the given available size
	// and returns the constrained size.
	Apply(available int) int
}

// Fixed is a constraint that represents a fixed size.
// The size will not grow or shrink based on available space.
type Fixed int

// Apply returns the fixed size, clamped to the available space.
func (f Fixed) Apply(available int) int {
	if f < 0 {
		return 0
	}
	if int(f) > available {
		return available
	}
	return int(f)
}

// Percent is a constraint that represents a percentage of the available size.
// Value should be between 0 and 100.
type Percent int

// Apply returns the percentage of the available size.
// Values outside 0-100 are clamped to the valid range.
func (p Percent) Apply(available int) int {
	if p <= 0 {
		return 0
	}
	if p >= 100 {
		return available
	}
	return available * int(p) / 100
}

// Ratio is a constraint that represents a ratio of the available size.
// For example, Ratio(1, 2) returns half the available space.
// Syntactic sugar for Percent(numerator * 100 / denominator).
func Ratio(numerator, denominator int) Percent {
	if denominator <= 0 {
		return 0
	}
	return Percent(numerator * 100 / denominator)
}

// Grow is a constraint that takes all remaining available space.
// When multiple Grow constraints are used, the space is divided equally.
type Grow struct{}

// Apply returns all available space.
func (g Grow) Apply(available int) int {
	return available
}

// MinSize is a constraint that applies a minimum size constraint.
// The size will be at least min pixels, but can grow if needed.
type MinSize struct {
	Min int
}

// Apply returns the minimum of min or available if available > 0.
func (m MinSize) Apply(available int) int {
	if available <= 0 {
		return m.Min
	}
	if m.Min > available {
		return available
	}
	return m.Min
}

// MaxSize is a constraint that applies a maximum size constraint.
// The size will be at most max pixels, but can shrink if needed.
type MaxSize struct {
	Max int
}

// Apply returns the minimum of max or available.
func (m MaxSize) Apply(available int) int {
	if available <= 0 {
		return 0
	}
	if m.Max <= 0 {
		return available
	}
	if m.Max < available {
		return m.Max
	}
	return available
}
