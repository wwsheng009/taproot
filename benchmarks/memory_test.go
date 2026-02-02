package benchmarks

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/wwsheng009/taproot/ui/components/messages"
	"github.com/wwsheng009/taproot/ui/list"
)

// BenchmarkMemory_ListItems measures memory allocation for list items
func BenchmarkMemory_ListItems(b *testing.B) {
	var m1, m2 runtime.MemStats

	b.ReportAllocs()
	runtime.GC()
	runtime.ReadMemStats(&m1)

	for i := 0; i < b.N; i++ {
		items := make([]list.FilterableItem, 1000)
		for j := 0; j < 1000; j++ {
			items[j] = list.NewListItem(
				fmt.Sprintf("id-%d", j),
				fmt.Sprintf("Item %d", j),
				fmt.Sprintf("Description %d", j),
			)
		}
		_ = items
	}

	runtime.GC()
	runtime.ReadMemStats(&m2)
	memoryUsed := m2.TotalAlloc - m1.TotalAlloc
	b.ReportMetric(float64(memoryUsed)/float64(b.N), "B/op_alloc")
}

// BenchmarkMemory_Messages measures memory allocation for messages
func BenchmarkMemory_Messages(b *testing.B) {
	var m1, m2 runtime.MemStats

	b.ReportAllocs()
	runtime.GC()
	runtime.ReadMemStats(&m1)

	for i := 0; i < b.N; i++ {
		msgs := make([]messages.MessageItem, 100)
		for j := 0; j < 100; j++ {
			if j%2 == 0 {
				msgs[j] = messages.NewUserMessage(
					fmt.Sprintf("user-%d", j),
					fmt.Sprintf("User message %d content here", j),
				)
			} else {
				msgs[j] = messages.NewAssistantMessage(
					fmt.Sprintf("assistant-%d", j),
					fmt.Sprintf("# Heading %d\n\nContent here", j),
				)
			}
		}
		_ = msgs
	}

	runtime.GC()
	runtime.ReadMemStats(&m2)
	memoryUsed := m2.TotalAlloc - m1.TotalAlloc
	b.ReportMetric(float64(memoryUsed)/float64(b.N), "B/op_alloc")
}

// BenchmarkMemory_FilteredResults measures memory for filter results
func BenchmarkMemory_FilteredResults(b *testing.B) {
	items := createTestItems(10000)
	filter := list.NewFilter()
	var m1, m2 runtime.MemStats

	b.ReportAllocs()
	runtime.GC()
	runtime.ReadMemStats(&m1)

	for i := 0; i < b.N; i++ {
		filter.SetQuery("item")
		results := filter.Apply(items)
		_ = results
	}

	runtime.GC()
	runtime.ReadMemStats(&m2)
	memoryUsed := m2.TotalAlloc - m1.TotalAlloc
	b.ReportMetric(float64(memoryUsed)/float64(b.N), "B/op_alloc")
}

// BenchmarkMemory_SelectedItems measures memory for selected items tracking
func BenchmarkMemory_SelectedItems(b *testing.B) {
	items := convertToItems(createTestItems(1000))
	selMgr := list.NewSelectionManager(list.SelectionModeMultiple)
	var m1, m2 runtime.MemStats

	b.ReportAllocs()
	runtime.GC()
	runtime.ReadMemStats(&m1)

	for i := 0; i < b.N; i++ {
		selMgr.SelectAll(items)
		ids := selMgr.SelectedIDs()
		_ = ids
		selMgr.Clear()
	}

	runtime.GC()
	runtime.ReadMemStats(&m2)
	memoryUsed := m2.TotalAlloc - m1.TotalAlloc
	b.ReportMetric(float64(memoryUsed)/float64(b.N), "B/op_alloc")
}

// BenchmarkMemory_GroupedItems measures memory for grouped items
func BenchmarkMemory_GroupedItems(b *testing.B) {
	groups := createTestGroups(100, 10)
	var m1, m2 runtime.MemStats

	b.ReportAllocs()
	runtime.GC()
	runtime.ReadMemStats(&m1)

	for i := 0; i < b.N; i++ {
		groupMgr := list.NewGroupManager()
		groupMgr.SetGroups(groups)
		_ = groupMgr
	}

	runtime.GC()
	runtime.ReadMemStats(&m2)
	memoryUsed := m2.TotalAlloc - m1.TotalAlloc
	b.ReportMetric(float64(memoryUsed)/float64(b.N), "B/op_alloc")
}

// BenchmarkMemory_RenderCache measures memory benefit from caching
func BenchmarkMemory_RenderCache(b *testing.B) {
	msg := messages.NewAssistantMessage("id-1", generateLargeContent(100))

	// First render (populate cache)
	_ = msg.View()

	var m1, m2 runtime.MemStats

	b.ReportAllocs()
	runtime.GC()
	runtime.ReadMemStats(&m1)

	// Subsequent renders (use cache)
	for i := 0; i < b.N; i++ {
		_ = msg.View()
	}

	runtime.GC()
	runtime.ReadMemStats(&m2)
	memoryUsed := m2.TotalAlloc - m1.TotalAlloc
	b.ReportMetric(float64(memoryUsed)/float64(b.N), "B/op_alloc")
}

// BenchmarkMemory_UpdateLoop measures memory in typical update loop
func BenchmarkMemory_UpdateLoop(b *testing.B) {
	items := convertToItems(createTestItems(100))
	filter := list.NewFilter()
	viewport := list.NewViewport(10, 100)
	selMgr := list.NewSelectionManager(list.SelectionModeMultiple)

	var m1, m2 runtime.MemStats

	b.ReportAllocs()
	runtime.GC()
	runtime.ReadMemStats(&m1)

	for i := 0; i < b.N; i++ {
		// Filter items
		filter.SetQuery(fmt.Sprintf("item-%d", i%10))
		filtered := filter.Apply(createTestItems(100))

		// Navigate viewport
		viewport.MoveDown()
		viewport.MoveUp()

		// Manage selection
		if i%20 == 0 {
			selMgr.SelectAll(items)
		} else {
			selMgr.Clear()
		}
		_ = filtered
	}

	runtime.GC()
	runtime.ReadMemStats(&m2)
	memoryUsed := m2.TotalAlloc - m1.TotalAlloc
	b.ReportMetric(float64(memoryUsed)/float64(b.N), "B/op_alloc")
}

// BenchmarkMemory_LargeMessageSet measures memory for large message sets
func BenchmarkMemory_LargeMessageSet(b *testing.B) {
	var m1, m2 runtime.MemStats

	b.ReportAllocs()
	runtime.GC()
	runtime.ReadMemStats(&m1)

	for i := 0; i < b.N; i++ {
		msgs := make([]messages.MessageItem, 500)
		for j := 0; j < 500; j++ {
			content := fmt.Sprintf("# Message %d\n\nContent: %s\n\n",
				j, generateRandomText(100))
			msgs[j] = messages.NewAssistantMessage(
				fmt.Sprintf("msg-%d", j),
				content,
			)
		}
		_ = msgs
	}

	runtime.GC()
	runtime.ReadMemStats(&m2)
	memoryUsed := m2.TotalAlloc - m1.TotalAlloc
	b.ReportMetric(float64(memoryUsed)/float64(b.N), "B/op_alloc")
}
