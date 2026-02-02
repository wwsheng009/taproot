package buffer

import (
	"testing"
)

// BenchmarkRenderWithOptimizations tests the optimized rendering
func BenchmarkRenderWithOptimizations(b *testing.B) {
	buf := GetBuffer(80, 24)
	defer PutBuffer(buf)

	// Fill buffer with content
	style := Style{Foreground: "202"}
	for y := 0; y < 24; y++ {
		for x := 0; x < 80; x++ {
			buf.cells[y][x] = Cell{
				Char:  'X',
				Width: 1,
				Style: style,
			}
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = buf.Render()
	}
}

// BenchmarkStyleCache tests style cache performance
func BenchmarkStyleCache(b *testing.B) {
	style := Style{Foreground: "202", Bold: true}

	cache := NewStyleCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cache.Get(style)
	}
}

// BenchmarkStyleCacheConcurrent tests concurrent style cache access
func BenchmarkStyleCacheConcurrent(b *testing.B) {
	cache := NewStyleCache()
	style := Style{Foreground: "202", Bold: true}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = cache.Get(style)
		}
	})
}

// BenchmarkLayoutManagerOptimized tests optimized layout manager
func BenchmarkLayoutManagerOptimized(b *testing.B) {
	lm := NewLayoutManager(80, 24)

	header := NewFillComponent('H', Style{Foreground: "202"})
	content := NewFillComponent('C', Style{})
	footer := NewFillComponent('F', Style{Foreground: "244"})

	lm.AddComponent("header", header)
	lm.AddComponent("content", content)
	lm.AddComponent("footer", footer)
	lm.CalculateLayout()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = lm.Render()
	}
}

// BenchmarkBufferPoolGetPut tests buffer pool performance
func BenchmarkBufferPoolGetPut(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := GetBuffer(80, 24)
			PutBuffer(buf)
		}
	})
}

// BenchmarkStringBuilderPool GetPut tests strings.Builder pool performance
func BenchmarkStringBuilderPoolGetPut(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sb := GetStringBuilder()
			sb.WriteString("test")
			PutStringBuilder(sb)
		}
	})
}

// BenchmarkRenderWithMixedStyles tests rendering with many different styles
func BenchmarkRenderWithMixedStyles(b *testing.B) {
	buf := GetBuffer(80, 24)
	defer PutBuffer(buf)

	// Fill buffer with mixed styles
	style1 := Style{Foreground: "202"}
	style2 := Style{Foreground: "201", Bold: true}
	style3 := Style{Foreground: "200", Italic: true}
	styles := []Style{style1, style2, style3, Style{}}

	for y := 0; y < 24; y++ {
		for x := 0; x < 80; x++ {
			buf.cells[y][x] = Cell{
				Char:  'X',
				Width: 1,
				Style: styles[(x+y)%4], // Mix styles across the buffer
			}
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = buf.Render()
	}
}
