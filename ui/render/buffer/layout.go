package buffer

// LayoutManager manages the layout of components using buffers
type LayoutManager struct {
	width      int
	height     int
	mainBuf    *Buffer
	layouts    map[string]*Rect
	components map[string]Renderable
}

// NewLayoutManager creates a new layout manager
func NewLayoutManager(width, height int) *LayoutManager {
	return &LayoutManager{
		width:      width,
		height:     height,
		mainBuf:    GetBuffer(width, height),
		layouts:    make(map[string]*Rect),
		components: make(map[string]Renderable),
	}
}

// SetSize sets the overall size of the layout
func (lm *LayoutManager) SetSize(width, height int) {
	if width != lm.width || height != lm.height {
		lm.width = width
		lm.height = height
		lm.mainBuf = NewBuffer(width, height)
	}
}

// AddComponent adds a component to the layout
func (lm *LayoutManager) AddComponent(id string, component Renderable) {
	lm.components[id] = component
}

// SetLayout sets the layout rectangle for a component
func (lm *LayoutManager) SetLayout(id string, rect Rect) {
	lm.layouts[id] = &rect
}

// CalculateLayout calculates and assigns layouts for all components
func (lm *LayoutManager) CalculateLayout() {
	// Header: Fixed at top, full width
	headerHeight := 5
	lm.SetLayout("header", Rect{
		X:      0,
		Y:      0,
		Width:  lm.width,
		Height: headerHeight,
	})

	// Footer: Fixed at bottom, full width
	footerHeight := 1
	lm.SetLayout("footer", Rect{
		X:      0,
		Y:      lm.height - footerHeight,
		Width:  lm.width,
		Height: footerHeight,
	})

	// Content: Middle area, full width, remaining height
	contentRect := Rect{
		X:      0,
		Y:      headerHeight,
		Width:  lm.width,
		Height: lm.height - headerHeight - footerHeight,
	}

	if contentRect.Height <= 0 {
		contentRect.Height = 0
	}
	lm.SetLayout("content", contentRect)
}

// ImageLayout calculates layout specifically for image viewer
func (lm *LayoutManager) ImageLayout(contentHeightHint int) {
	// Header: Fixed at top
	headerHeight := 5
	lm.SetLayout("header", Rect{
		X:      0,
		Y:      0,
		Width:  lm.width,
		Height: headerHeight,
	})

	// Footer: Fixed at bottom
	footerHeight := 1
	lm.SetLayout("footer", Rect{
		X:      0,
		Y:      lm.height - footerHeight,
		Width:  lm.width,
		Height: footerHeight,
	})

	// Content: Middle area
	totalContentSpace := lm.height - headerHeight - footerHeight
	if totalContentSpace < 0 {
		totalContentSpace = 0
	}

	contentHeight := totalContentSpace
	if contentHeightHint > 0 && contentHeightHint < totalContentSpace {
		contentHeight = contentHeightHint
	}

	// Center vertical content if smaller than available space
	contentY := headerHeight
	if contentHeight < totalContentSpace {
		paddingTop := (totalContentSpace - contentHeight) / 2
		contentY += paddingTop
	}

	lm.SetLayout("content", Rect{
		X:      0,
		Width:  lm.width,
		Y:      contentY,
		Height: contentHeight,
	})
}

// Render renders all components to their assigned layouts
func (lm *LayoutManager) Render() string {
	// Clear the main buffer
	lm.mainBuf = GetBuffer(lm.width, lm.height)

	// Render each component to its assigned area
	for id, component := range lm.components {
		layoutRect, exists := lm.layouts[id]
		if !exists {
			continue
		}

		// Create a sub-buffer for this component
		componentBuf := GetBuffer(layoutRect.Width, layoutRect.Height)

		// Let the component render to its sub-buffer
		component.Render(componentBuf, Rect{
			X:      0,
			Y:      0,
			Width:  layoutRect.Width,
			Height: layoutRect.Height,
		})

		// Write the component buffer to the main buffer at the correct position
		lm.mainBuf.WriteBuffer(Point{
			X: layoutRect.X,
			Y: layoutRect.Y,
		}, componentBuf)

		PutBuffer(componentBuf)
	}

	result := lm.mainBuf.Render()
	PutBuffer(lm.mainBuf)

	return result
}

// GetBuffer returns the main buffer (useful for debugging)
func (lm *LayoutManager) GetBuffer() *Buffer {
	return lm.mainBuf
}
