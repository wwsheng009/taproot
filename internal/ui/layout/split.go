package layout

import "image"

// SplitVertical splits the area vertically into top and bottom parts.
// The split is based on the provided constraint.
//
// Example:
//
//	top, bottom := SplitVertical(area, Fixed(5))     // 5 units high at top
//	top, bottom := SplitVertical(area, Ratio(1, 2))  // Top half, bottom half
//	top, bottom := SplitVertical(area, Percent(30))  // Top 30%, bottom 70%
func SplitVertical(area Area, constraint Constraint) (top, bottom Area) {
	height := min(constraint.Apply(area.Dy()), area.Dy())

	top = Area{
		Min: area.Rect().Min,
		Max: image.Pt(area.Rect().Max.X, area.Rect().Min.Y+height),
	}

	bottom = Area{
		Min: image.Pt(area.Rect().Min.X, area.Rect().Min.Y+height),
		Max: area.Rect().Max,
	}

	return top, bottom
}

// SplitHorizontal splits the area horizontally into left and right parts.
// The split is based on the provided constraint.
//
// Example:
//
//	left, right := SplitHorizontal(area, Fixed(30))    // 30 units wide on left
//	left, right := SplitHorizontal(area, Ratio(1, 3))  // Left 1/3, right 2/3
//	left, right := SplitHorizontal(area, Percent(25))  // Left 25%, right 75%
func SplitHorizontal(area Area, constraint Constraint) (left, right Area) {
	width := min(constraint.Apply(area.Dx()), area.Dx())

	left = Area{
		Min: area.Rect().Min,
		Max: image.Pt(area.Rect().Min.X+width, area.Rect().Max.Y),
	}

	right = Area{
		Min: image.Pt(area.Rect().Min.X+width, area.Rect().Min.Y),
		Max: area.Rect().Max,
	}

	return left, right
}

// CenterRect returns a new Area centered within the given area
// with the specified width and height.
func CenterRect(area Area, width, height int) Area {
	centerX := area.Rect().Min.X + area.Dx()/2
	centerY := area.Rect().Min.Y + area.Dy()/2

	minX := centerX - width/2
	minY := centerY - height/2
	maxX := minX + width
	maxY := minY + height

	return Area(image.Rect(minX, minY, maxX, maxY)).Intersect(area)
}

// TopLeftRect returns a new Area positioned at the top-left corner
// of the given area with the specified width and height.
func TopLeftRect(area Area, width, height int) Area {
	return Area(image.Rect(area.Rect().Min.X, area.Rect().Min.Y, area.Rect().Min.X+width, area.Rect().Min.Y+height)).Intersect(area)
}

// TopCenterRect returns a new Area positioned at the top-center
// of the given area with the specified width and height.
func TopCenterRect(area Area, width, height int) Area {
	centerX := area.Rect().Min.X + area.Dx()/2
	minX := centerX - width/2

	return Area(image.Rect(minX, area.Rect().Min.Y, minX+width, area.Rect().Min.Y+height)).Intersect(area)
}

// TopRightRect returns a new Area positioned at the top-right corner
// of the given area with the specified width and height.
func TopRightRect(area Area, width, height int) Area {
	return Area(image.Rect(area.Rect().Max.X-width, area.Rect().Min.Y, area.Rect().Max.X, area.Rect().Min.Y+height)).Intersect(area)
}

// RightCenterRect returns a new Area positioned at the right-center
// of the given area with the specified width and height.
func RightCenterRect(area Area, width, height int) Area {
	centerY := area.Rect().Min.Y + area.Dy()/2
	minY := centerY - height/2

	return Area(image.Rect(area.Rect().Max.X-width, minY, area.Rect().Max.X, minY+height)).Intersect(area)
}

// LeftCenterRect returns a new Area positioned at the left-center
// of the given area with the specified width and height.
func LeftCenterRect(area Area, width, height int) Area {
	centerY := area.Rect().Min.Y + area.Dy()/2
	minY := centerY - height/2

	return Area(image.Rect(area.Rect().Min.X, minY, area.Rect().Min.X+width, minY+height)).Intersect(area)
}

// BottomLeftRect returns a new Area positioned at the bottom-left corner
// of the given area with the specified width and height.
func BottomLeftRect(area Area, width, height int) Area {
	return Area(image.Rect(area.Rect().Min.X, area.Rect().Max.Y-height, area.Rect().Min.X+width, area.Rect().Max.Y)).Intersect(area)
}

// BottomCenterRect returns a new Area positioned at the bottom-center
// of the given area with the specified width and height.
func BottomCenterRect(area Area, width, height int) Area {
	centerX := area.Rect().Min.X + area.Dx()/2
	minX := centerX - width/2

	return Area(image.Rect(minX, area.Rect().Max.Y-height, minX+width, area.Rect().Max.Y)).Intersect(area)
}

// BottomRightRect returns a new Area positioned at the bottom-right corner
// of the given area with the specified width and height.
func BottomRightRect(area Area, width, height int) Area {
	return Area(image.Rect(area.Rect().Max.X-width, area.Rect().Max.Y-height, area.Rect().Max.X, area.Rect().Max.Y)).Intersect(area)
}

// Pad returns a new Area with padding applied to all sides.
// The padding value is applied to all four sides.
func Pad(area Area, padding int) Area {
	rect := area.Rect()
	return Area(image.Rect(
		rect.Min.X+padding,
		rect.Min.Y+padding,
		rect.Max.X-padding,
		rect.Max.Y-padding,
	)).Intersect(area)
}

// Inset returns a new Area with different padding for each side.
func Inset(area Area, top, right, bottom, left int) Area {
	rect := area.Rect()
	return Area(image.Rect(
		rect.Min.X+left,
		rect.Min.Y+top,
		rect.Max.X-right,
		rect.Max.Y-bottom,
	)).Intersect(area)
}

// WithPadding is a functional option for applying padding.
type WithPadding func(Area) Area

// Padding returns a WithPadding function that applies uniform padding.
func Padding(amount int) WithPadding {
	return func(area Area) Area {
		return Pad(area, amount)
	}
}

// HorizontalPadding returns a WithPadding function that applies padding
// to left and right sides only.
func HorizontalPadding(amount int) WithPadding {
	return func(area Area) Area {
		return Inset(area, 0, amount, 0, amount)
	}
}

// VerticalPadding returns a WithPadding function that applies padding
// to top and bottom sides only.
func VerticalPadding(amount int) WithPadding {
	return func(area Area) Area {
		return Inset(area, amount, 0, amount, 0)
	}
}
