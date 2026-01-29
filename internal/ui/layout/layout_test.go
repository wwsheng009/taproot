package layout

import (
	"testing"
)

func TestNewArea(t *testing.T) {
	area := NewArea(0, 0, 10, 10)
	if area.Dx() != 10 {
		t.Errorf("expected width 10, got %d", area.Dx())
	}
	if area.Dy() != 10 {
		t.Errorf("expected height 10, got %d", area.Dy())
	}
	if area.Empty() {
		t.Error("area should not be empty")
	}
}

func TestAreaTopLeft(t *testing.T) {
	area := NewArea(5, 3, 15, 13)
	topLeft := area.TopLeft()
	if topLeft.X != 5 {
		t.Errorf("expected X 5, got %d", topLeft.X)
	}
	if topLeft.Y != 3 {
		t.Errorf("expected Y 3, got %d", topLeft.Y)
	}
}

func TestAreaBottomRight(t *testing.T) {
	area := NewArea(5, 3, 15, 13)
	bottomRight := area.BottomRight()
	if bottomRight.X != 15 {
		t.Errorf("expected X 15, got %d", bottomRight.X)
	}
	if bottomRight.Y != 13 {
		t.Errorf("expected Y 13, got %d", bottomRight.Y)
	}
}

func TestAreaEmpty(t *testing.T) {
	tests := []struct {
		name string
		area Area
		empty bool
	}{
		{"normal area", NewArea(0, 0, 10, 10), false},
		{"zero width", NewArea(0, 0, 0, 10), true},
		{"zero height", NewArea(0, 0, 10, 0), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.area.Empty(); got != tt.empty {
				t.Errorf("Empty() = %v, want %v", got, tt.empty)
			}
		})
	}
}

func TestFixedConstraint(t *testing.T) {
	f := Fixed(10)
	tests := []struct {
		available int
		expected   int
	}{
		{20, 10},
		{10, 10},
		{5, 5}, // clamped to available
	}
	for _, tt := range tests {
		if got := f.Apply(tt.available); got != tt.expected {
			t.Errorf("Apply(%d) = %d, want %d", tt.available, got, tt.expected)
		}
	}
}

func TestPercentConstraint(t *testing.T) {
	tests := []struct {
		percent   Percent
		available int
		expected  int
	}{
		{50, 100, 50},
		{25, 80, 20},
		{100, 100, 100},
		{0, 100, 0},
	}
	for _, tt := range tests {
		if got := tt.percent.Apply(tt.available); got != tt.expected {
			t.Errorf("Percent(%d).Apply(%d) = %d, want %d", tt.percent, tt.available, got, tt.expected)
		}
	}
}

func TestRatio(t *testing.T) {
	tests := []struct {
		numerator   int
		denominator int
		available   int
		expected    int
	}{
		{1, 2, 100, 50},  // half
		{1, 3, 90, 30},   // ~1/3 (with integer division)
		{2, 3, 90, 60},   // ~2/3
		{1, 4, 100, 25},  // quarter
	}
	for _, tt := range tests {
		ratio := Ratio(tt.numerator, tt.denominator)
		got := ratio.Apply(tt.available)
		if got < tt.expected-1 || got > tt.expected+1 {
			t.Errorf("Ratio(%d, %d).Apply(%d) = %d, want %d (±1)", tt.numerator, tt.denominator, tt.available, got, tt.expected)
		}
	}
}

func TestGrowConstraint(t *testing.T) {
	g := Grow{}
	tests := []struct {
		available int
		expected  int
	}{
		{100, 100},
		{50, 50},
		{0, 0},
	}
	for _, tt := range tests {
		if got := g.Apply(tt.available); got != tt.expected {
			t.Errorf("Grow{}.Apply(%d) = %d, want %d", tt.available, got, tt.expected)
		}
	}
}

func TestSplitVertical(t *testing.T) {
	area := NewArea(0, 0, 20, 20)

	top, bottom := SplitVertical(area, Fixed(10))
	if top.Dx() != 20 {
		t.Errorf("top should have width 20, got %d", top.Dx())
	}
	if top.Dy() != 10 {
		t.Errorf("top should have height 10, got %d", top.Dy())
	}
	if bottom.Dx() != 20 {
		t.Errorf("bottom should have width 20, got %d", bottom.Dx())
	}
	if bottom.Dy() != 10 {
		t.Errorf("bottom should have height 10, got %d", bottom.Dy())
	}
}

func TestSplitHorizontal(t *testing.T) {
	area := NewArea(0, 0, 20, 20)

	left, right := SplitHorizontal(area, Fixed(10))
	if left.Dx() != 10 {
		t.Errorf("left should have width 10, got %d", left.Dx())
	}
	if left.Dy() != 20 {
		t.Errorf("left should have height 20, got %d", left.Dy())
	}
	if right.Dx() != 10 {
		t.Errorf("right should have width 10, got %d", right.Dx())
	}
	if right.Dy() != 20 {
		t.Errorf("right should have height 20, got %d", right.Dy())
	}
}

func TestSplitVerticalPercent(t *testing.T) {
	area := NewArea(0, 0, 100, 100)

	top, bottom := SplitVertical(area, Percent(50))
	if top.Dy() != 50 {
		t.Errorf("top height: expected 50, got %d", top.Dy())
	}
	if bottom.Dy() != 50 {
		t.Errorf("bottom height: expected 50, got %d", bottom.Dy())
	}
}

func TestSplitHorizontalPercent(t *testing.T) {
	area := NewArea(0, 0, 100, 100)

	left, right := SplitHorizontal(area, Percent(25))
	if left.Dx() != 25 {
		t.Errorf("left width: expected 25, got %d", left.Dx())
	}
	if right.Dx() != 75 {
		t.Errorf("right width: expected 75, got %d", right.Dx())
	}
}

func TestCenterRect(t *testing.T) {
	area := NewArea(0, 0, 100, 100)
	centered := CenterRect(area, 20, 20)

	if centered.Dx() != 20 {
		t.Errorf("expected width 20, got %d", centered.Dx())
	}
	if centered.Dy() != 20 {
		t.Errorf("expected height 20, got %d", centered.Dy())
	}
	// Check centering
	topLeft := centered.TopLeft()
	if topLeft.X != 40 {
		t.Errorf("expected X 40, got %d", topLeft.X)
	}
	if topLeft.Y != 40 {
		t.Errorf("expected Y 40, got %d", topLeft.Y)
	}
}

func TestPad(t *testing.T) {
	area := NewArea(0, 0, 100, 100)
	padded := Pad(area, 10)

	if padded.Dx() != 80 {
		t.Errorf("expected width 80, got %d", padded.Dx())
	}
	if padded.Dy() != 80 {
		t.Errorf("expected height 80, got %d", padded.Dy())
	}

	topLeft := padded.TopLeft()
	if topLeft.X != 10 {
		t.Errorf("expected X 10, got %d", topLeft.X)
	}
	if topLeft.Y != 10 {
		t.Errorf("expected Y 10, got %d", topLeft.Y)
	}
}

func TestInset(t *testing.T) {
	area := NewArea(0, 0, 100, 100)
	inset := Inset(area, 5, 10, 15, 20)

	if inset.Dx() != 70 {
		t.Errorf("expected width 70, got %d", inset.Dx())
	}
	if inset.Dy() != 80 {
		t.Errorf("expected height 80, got %d", inset.Dy())
	}

	topLeft := inset.TopLeft()
	if topLeft.X != 20 {
		t.Errorf("expected X 20, got %d", topLeft.X)
	}
	if topLeft.Y != 5 {
		t.Errorf("expected Y 5, got %d", topLeft.Y)
	}
}

func TestRowLayout(t *testing.T) {
	area := NewArea(0, 0, 100, 20)

	children := []FlexChild{
		NewFlexChild(Fixed(20)),
		NewFlexChild(Fixed(30)),
	}

	areas := RowLayout(area, children)

	if len(areas) != 2 {
		t.Fatalf("expected 2 areas, got %d", len(areas))
	}
	if areas[0].Dx() != 20 {
		t.Errorf("first area should have width 20, got %d", areas[0].Dx())
	}
	if areas[1].Dx() != 30 {
		t.Errorf("second area should have width 30, got %d", areas[1].Dx())
	}
}

func TestRowLayoutWithGrow(t *testing.T) {
	area := NewArea(0, 0, 100, 20)

	children := []FlexChild{
		NewFlexChild(Fixed(20)),
		NewFlexChild(Grow{}).WithGrow(),
	}

	areas := RowLayout(area, children)

	if len(areas) != 2 {
		t.Fatalf("expected 2 areas, got %d", len(areas))
	}
	if areas[0].Dx() != 20 {
		t.Errorf("first area should have width 20, got %d", areas[0].Dx())
	}
	if areas[1].Dx() != 80 {
		t.Errorf("second area should have width 80, got %d", areas[1].Dx())
	}
}

func TestColumnLayout(t *testing.T) {
	area := NewArea(0, 0, 20, 100)

	children := []FlexChild{
		NewFlexChild(Fixed(10)),
		NewFlexChild(Fixed(20)),
	}

	areas := ColumnLayout(area, children)

	if len(areas) != 2 {
		t.Fatalf("expected 2 areas, got %d", len(areas))
	}
	if areas[0].Dy() != 10 {
		t.Errorf("first area should have height 10, got %d", areas[0].Dy())
	}
	if areas[1].Dy() != 20 {
		t.Errorf("second area should have height 20, got %d", areas[1].Dy())
	}
}

func TestColumnLayoutWithGrow(t *testing.T) {
	area := NewArea(0, 0, 20, 100)

	children := []FlexChild{
		NewFlexChild(Fixed(10)),
		NewFlexChild(Grow{}).WithGrow(),
	}

	areas := ColumnLayout(area, children)

	if len(areas) != 2 {
		t.Fatalf("expected 2 areas, got %d", len(areas))
	}
	if areas[0].Dy() != 10 {
		t.Errorf("first area should have height 10, got %d", areas[0].Dy())
	}
	if areas[1].Dy() != 90 {
		t.Errorf("second area should have height 90, got %d", areas[1].Dy())
	}
}

func TestFlexRow(t *testing.T) {
	area := NewArea(0, 0, 100, 20)
	areas := FlexRow(area, Fixed(20), Fixed(30), Fixed(50))

	if len(areas) != 3 {
		t.Fatalf("expected 3 areas, got %d", len(areas))
	}
	if areas[0].Dx() != 20 {
		t.Errorf("first area should have width 20, got %d", areas[0].Dx())
	}
	if areas[1].Dx() != 30 {
		t.Errorf("second area should have width 30, got %d", areas[1].Dx())
	}
	if areas[2].Dx() != 50 {
		t.Errorf("third area should have width 50, got %d", areas[2].Dx())
	}
}

func TestFlexColumn(t *testing.T) {
	area := NewArea(0, 0, 20, 100)
	areas := FlexColumn(area, Fixed(20), Fixed(30), Fixed(50))

	if len(areas) != 3 {
		t.Fatalf("expected 3 areas, got %d", len(areas))
	}
	if areas[0].Dy() != 20 {
		t.Errorf("first area should have height 20, got %d", areas[0].Dy())
	}
	if areas[1].Dy() != 30 {
		t.Errorf("second area should have height 30, got %d", areas[1].Dy())
	}
	if areas[2].Dy() != 50 {
		t.Errorf("third area should have height 50, got %d", areas[2].Dy())
	}
}

func TestGridLayout(t *testing.T) {
	area := NewArea(0, 0, 60, 40)
	config := NewGridConfig(2, 3)

	cells := GridLayout(area, config)

	if len(cells) != 6 {
		t.Fatalf("expected 6 cells, got %d", len(cells))
	}

	// Check first cell (row 0, col 0)
	if cells[0].Dx() != 20 {
		t.Errorf("cell[0] width: expected 20, got %d", cells[0].Dx())
	}
	if cells[0].Dy() != 20 {
		t.Errorf("cell[0] height: expected 20, got %d", cells[0].Dy())
	}

	// Check last cell (row 1, col 2)
	if cells[5].Dx() != 20 {
		t.Errorf("cell[5] width: expected 20, got %d", cells[5].Dx())
	}
	if cells[5].Dy() != 20 {
		t.Errorf("cell[5] height: expected 20, got %d", cells[5].Dy())
	}
}

func TestGridLayoutWithGaps(t *testing.T) {
	area := NewArea(0, 0, 60, 40)
	config := NewGridConfig(2, 3).WithColGaps(1).WithRowGaps(1)

	cells := GridLayout(area, config)

	if len(cells) != 6 {
		t.Fatalf("expected 6 cells, got %d", len(cells))
	}

	// With gaps, cells should be smaller
	// Available: 60 width, 40 height
	// Gaps: 2 col gaps (2), 1 row gap (1)
	// Available for cells: 58 width, 39 height
	// Cell size: ~19 x ~19
	if cells[0].Dx() < 19 || cells[0].Dx() > 20 {
		t.Errorf("cell[0] width should be ~19, got %d", cells[0].Dx())
	}
	if cells[0].Dy() < 19 || cells[0].Dy() > 20 {
		t.Errorf("cell[0] height should be ~19, got %d", cells[0].Dy())
	}
}

func TestGetCell(t *testing.T) {
	area := NewArea(0, 0, 60, 40)
	config := NewGridConfig(2, 3)
	cells := GridLayout(area, config)

	// Valid cell
	cell := GetCell(cells, config, 0, 0)
	if cell.Empty() {
		t.Error("cell (0,0) should not be empty")
	}

	// Invalid row
	cell = GetCell(cells, config, 5, 0)
	if !cell.Empty() {
		t.Error("cell (5,0) should be empty")
	}

	// Invalid column
	cell = GetCell(cells, config, 0, 5)
	if !cell.Empty() {
		t.Error("cell (0,5) should be empty")
	}
}

func TestGetRow(t *testing.T) {
	area := NewArea(0, 0, 60, 40)
	config := NewGridConfig(2, 3)
	cells := GridLayout(area, config)

	row := GetRow(cells, config, 0)
	if len(row) != 3 {
		t.Errorf("expected 3 cells in row 0, got %d", len(row))
	}

	// Invalid row
	row = GetRow(cells, config, 5)
	if len(row) != 0 {
		t.Errorf("expected 0 cells in row 5, got %d", len(row))
	}
}

func TestGetColumn(t *testing.T) {
	area := NewArea(0, 0, 60, 40)
	config := NewGridConfig(2, 3)
	cells := GridLayout(area, config)

	col := GetColumn(cells, config, 0)
	if len(col) != 2 {
		t.Errorf("expected 2 cells in column 0, got %d", len(col))
	}

	// Invalid column
	col = GetColumn(cells, config, 5)
	if len(col) != 0 {
		t.Errorf("expected 0 cells in column 5, got %d", len(col))
	}
}

func TestSpanCell(t *testing.T) {
	area := NewArea(0, 0, 60, 40)
	config := NewGridConfig(2, 3)
	cells := GridLayout(area, config)

	// Span 2 cells horizontally starting from (0, 0)
	span := SpanCell(cells, config, 0, 0, 2, 1)

	// Span should cover 2 columns (40 width) and 1 row (20 height)
	if span.Dx() < 39 || span.Dx() > 41 {
		t.Errorf("span width: expected ~40, got %d", span.Dx())
	}
	if span.Dy() < 19 || span.Dy() > 21 {
		t.Errorf("span height: expected ~20, got %d", span.Dy())
	}

	// Span 2 cells vertically starting from (0, 0)
	span = SpanCell(cells, config, 0, 0, 1, 2)

	// Span should cover 1 column (20 width) and 2 rows (40 height)
	if span.Dx() < 19 || span.Dx() > 21 {
		t.Errorf("span width: expected ~20, got %d", span.Dx())
	}
	if span.Dy() < 39 || span.Dy() > 41 {
		t.Errorf("span height: expected ~40, got %d", span.Dy())
	}
}

func TestFixedGrid(t *testing.T) {
	area := NewArea(0, 0, 100, 100)
	cells := FixedGrid(area, 20, 20)

	// Should fit 5x5 cells
	expectedCells := 25
	if len(cells) != expectedCells {
		t.Errorf("expected %d cells, got %d", expectedCells, len(cells))
	}

	for _, cell := range cells {
		if cell.Dx() != 20 {
			t.Errorf("cell width: expected 20, got %d", cell.Dx())
		}
		if cell.Dy() != 20 {
			t.Errorf("cell height: expected 20, got %d", cell.Dy())
		}
	}
}

func TestUniformGrid(t *testing.T) {
	area := NewArea(0, 0, 60, 40)
	cells := UniformGrid(area, 2, 3)

	if len(cells) != 6 {
		t.Fatalf("expected 6 cells, got %d", len(cells))
	}

	for _, cell := range cells {
		if cell.Dx() != 20 {
			t.Errorf("cell width: expected 20, got %d", cell.Dx())
		}
		if cell.Dy() != 20 {
			t.Errorf("cell height: expected 20, got %d", cell.Dy())
		}
	}
}

func TestEmptyInputs(t *testing.T) {
	t.Run("empty grid with zero rows", func(t *testing.T) {
		config := NewGridConfig(0, 3)
		cells := GridLayout(NewArea(0, 0, 100, 100), config)
		if len(cells) != 0 {
			t.Errorf("expected 0 cells, got %d", len(cells))
		}
	})

	t.Run("empty grid with zero cols", func(t *testing.T) {
		config := NewGridConfig(2, 0)
		cells := GridLayout(NewArea(0, 0, 100, 100), config)
		if len(cells) != 0 {
			t.Errorf("expected 0 cells, got %d", len(cells))
		}
	})

	t.Run("empty area", func(t *testing.T) {
		cells := GridLayout(NewArea(0, 0, 0, 0), NewGridConfig(2, 2))
		if len(cells) != 0 {
			t.Errorf("expected 0 cells, got %d", len(cells))
		}
	})

	t.Run("empty flex", func(t *testing.T) {
		area := NewArea(0, 0, 100, 20)
		areas := RowLayout(area, []FlexChild{})
		if len(areas) != 0 {
			t.Errorf("expected 0 areas, got %d", len(areas))
		}
	})
}
