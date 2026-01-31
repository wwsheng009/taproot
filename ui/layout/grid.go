package layout

import "image"

// GridConfig defines the configuration for a grid layout.
type GridConfig struct {
	// Rows is the number of rows in the grid.
	Rows int
	// Cols is the number of columns in the grid.
	Cols int
	// RowGaps is the number of empty rows between each row.
	RowGaps int
	// ColGaps is the number of empty columns between each column.
	ColGaps int
}

// NewGridConfig creates a new GridConfig with the specified rows and columns.
func NewGridConfig(rows, cols int) GridConfig {
	return GridConfig{
		Rows:    rows,
		Cols:    cols,
		RowGaps: 0,
		ColGaps: 0,
	}
}

// WithRowGaps sets the row gap and returns the config.
func (c GridConfig) WithRowGaps(gap int) GridConfig {
	c.RowGaps = gap
	return c
}

// WithColGaps sets the column gap and returns the config.
func (c GridConfig) WithColGaps(gap int) GridConfig {
	c.ColGaps = gap
	return c
}

// GridLayout creates a grid layout within the specified area.
// The cells are returned in row-major order (left to right, top to bottom).
//
// Example:
//
//	config := layout.NewGridConfig(2, 3).WithColGaps(1)
//	cells := layout.GridLayout(area, config)
//	// cells[0] = row 0, col 0
//	// cells[1] = row 0, col 1
//	// cells[2] = row 0, col 2
//	// cells[3] = row 1, col 0
//	// cells[4] = row 1, col 1
//	// cells[5] = row 1, col 2
func GridLayout(area Area, config GridConfig) []Area {
	if config.Rows <= 0 || config.Cols <= 0 {
		return []Area{}
	}

	if area.Empty() {
		return []Area{}
	}

	// Calculate total gaps
	totalRowGaps := config.RowGaps * (config.Rows - 1)
	totalColGaps := config.ColGaps * (config.Cols - 1)

	// Calculate available space for cells
	availableWidth := area.Dx() - totalColGaps
	availableHeight := area.Dy() - totalRowGaps

	if availableWidth <= 0 || availableHeight <= 0 {
		return []Area{}
	}

	// Calculate cell sizes
	cellWidth := availableWidth / config.Cols
	cellHeight := availableHeight / config.Rows

	// Create cells
	totalCells := config.Rows * config.Cols
	cells := make([]Area, 0, totalCells)

	for row := 0; row < config.Rows; row++ {
		for col := 0; col < config.Cols; col++ {
			// Calculate top-left position
			x := area.Rect().Min.X + col*(cellWidth+config.ColGaps)
			y := area.Rect().Min.Y + row*(cellHeight+config.RowGaps)

			// Calculate bottom-right position
			maxX := x + cellWidth
			maxY := y + cellHeight

			// Create cell area
			cell := NewArea(x, y, maxX, maxY)

			// Ensure cell is within parent area
			cell = cell.Intersect(area)

			cells = append(cells, cell)
		}
	}

	return cells
}

// GetCell returns the cell at the specified row and column.
// Returns an empty Area if the coordinates are out of bounds.
func GetCell(cells []Area, config GridConfig, row, col int) Area {
	if row < 0 || row >= config.Rows || col < 0 || col >= config.Cols {
		return NewArea(0, 0, 0, 0)
	}

	index := row*config.Cols + col
	if index >= len(cells) {
		return NewArea(0, 0, 0, 0)
	}

	return cells[index]
}

// GetRow returns all cells in the specified row.
func GetRow(cells []Area, config GridConfig, row int) []Area {
	if row < 0 || row >= config.Rows {
		return []Area{}
	}

	start := row * config.Cols
	end := start + config.Cols
	if end > len(cells) {
		end = len(cells)
	}

	return cells[start:end]
}

// GetColumn returns all cells in the specified column.
func GetColumn(cells []Area, config GridConfig, col int) []Area {
	if col < 0 || col >= config.Cols {
		return []Area{}
	}

	cols := make([]Area, 0, config.Rows)
	for row := 0; row < config.Rows; row++ {
		index := row*config.Cols + col
		if index < len(cells) {
			cols = append(cols, cells[index])
		}
	}

	return cols
}

// SpanCell returns an area that spans multiple cells starting from the specified position.
// The span width and height are specified in cells, not pixels.
func SpanCell(cells []Area, config GridConfig, startRow, startCol, spanCols, spanRows int) Area {
	if startRow < 0 || startCol < 0 || startRow >= config.Rows || startCol >= config.Cols {
		return NewArea(0, 0, 0, 0)
	}

	if spanCols <= 0 || spanRows <= 0 {
		return NewArea(0, 0, 0, 0)
	}

	startCell := GetCell(cells, config, startRow, startCol)
	if startCell.Empty() {
		return NewArea(0, 0, 0, 0)
	}

	// Calculate the outer bounds of the span
	endRow := min(startRow+spanRows-1, config.Rows-1)
	endCol := min(startCol+spanCols-1, config.Cols-1)

	endCell := GetCell(cells, config, endRow, endCol)
	if endCell.Empty() {
		return NewArea(0, 0, 0, 0)
	}

	// Combine the areas
	startPt := image.Pt(startCell.TopLeft().X, startCell.TopLeft().Y)
	endPt := image.Pt(endCell.BottomRight().X, endCell.BottomRight().Y)
	span := NewArea(startPt.X, startPt.Y, endPt.X, endPt.Y)

	return span
}

// FixedGrid creates a grid with fixed cell sizes.
// This is useful when you want cells of a specific size rather than filling the area.
func FixedGrid(area Area, cellWidth, cellHeight int) []Area {
	if area.Empty() || cellWidth <= 0 || cellHeight <= 0 {
		return []Area{}
	}

	rows := area.Dy() / cellHeight
	cols := area.Dx() / cellWidth
	if rows <= 0 || cols <= 0 {
		return []Area{}
	}

	cells := make([]Area, 0, rows*cols)

	for y := area.Rect().Min.Y; y+cellHeight <= area.Rect().Max.Y; y += cellHeight {
		for x := area.Rect().Min.X; x+cellWidth <= area.Rect().Max.X; x += cellWidth {
			cell := NewArea(x, y, x+cellWidth, y+cellHeight)
			cells = append(cells, cell)
		}
	}

	return cells
}

// UniformGrid creates a grid with all cells having the same size.
// Cells are distributed evenly across the area.
func UniformGrid(area Area, rows, cols int) []Area {
	return GridLayout(area, NewGridConfig(rows, cols))
}
