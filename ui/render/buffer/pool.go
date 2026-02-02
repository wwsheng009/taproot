package buffer

import (
	"strings"
	"sync"
)

// BufferPool manages a pool of reusable Buffer objects
var bufferPool = sync.Pool{
	New: func() interface{} {
		return &Buffer{
			cells: make([][]Cell, 0, 24),
		}
	},
}

// StringBuilderPool manages a pool of strings.Builder objects
var stringBuilderPool = sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
}

// GetBuffer retrieves a buffer from the pool
func GetBuffer(width, height int) *Buffer {
	buf := bufferPool.Get().(*Buffer)

	// Reuse or resize if needed
	if buf.width == width && buf.height == height && len(buf.cells) == height {
		// Clear the buffer
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				buf.cells[y][x] = Cell{Char: ' ', Width: 1}
			}
		}
	} else {
		// Initialize buffer
		buf.width = width
		buf.height = height

		// Resize cells slice
		if cap(buf.cells) >= height {
			buf.cells = buf.cells[:height]
		} else {
			buf.cells = make([][]Cell, height)
		}

		// Initialize each row
		for y := 0; y < height; y++ {
			if len(buf.cells[y]) < width || cap(buf.cells[y]) < width {
				buf.cells[y] = make([]Cell, width)
			}
			for x := 0; x < width; x++ {
				buf.cells[y][x] = Cell{Char: ' ', Width: 1}
			}
		}
	}

	return buf
}

// PutBuffer returns a buffer to the pool
func PutBuffer(buf *Buffer) {
	if buf != nil {
		bufferPool.Put(buf)
	}
}

// GetStringBuilder retrieves a strings.Builder from the pool
func GetStringBuilder() *strings.Builder {
	sb := stringBuilderPool.Get().(*strings.Builder)
	sb.Reset()
	return sb
}

// PutStringBuilder returns a strings.Builder to the pool
func PutStringBuilder(sb *strings.Builder) {
	if sb != nil {
		stringBuilderPool.Put(sb)
	}
}
