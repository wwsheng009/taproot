package benchmarks

import (
	"fmt"
	"testing"

	"github.com/wwsheng009/taproot/ui/list"
)

// BenchmarkListItemCreation measures performance of creating list items
func BenchmarkListItemCreation(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = list.NewListItem("id-1", "Item Title", "Item Description")
	}
}

// BenchmarkListFilter_Small measures filtering performance with 100 items
func BenchmarkListFilter_Small(b *testing.B) {
	items := createTestItems(100)
	filter := list.NewFilter()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		filter.SetQuery("item")
		_ = filter.Apply(items)
	}
}

// BenchmarkListFilter_Medium measures filtering performance with 1000 items
func BenchmarkListFilter_Medium(b *testing.B) {
	items := createTestItems(1000)
	filter := list.NewFilter()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		filter.SetQuery("item")
		_ = filter.Apply(items)
	}
}

// BenchmarkListFilter_Large measures filtering performance with 10000 items
func BenchmarkListFilter_Large(b *testing.B) {
	items := createTestItems(10000)
	filter := list.NewFilter()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		filter.SetQuery("item")
		_ = filter.Apply(items)
	}
}

// BenchmarkListSelection_BulkSelect measures bulk selection performance
func BenchmarkListSelection_BulkSelect(b *testing.B) {
	items := convertToItems(createTestItems(1000))
	selMgr := list.NewSelectionManager(list.SelectionModeMultiple)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		selMgr.SelectAll(items)
		selMgr.Clear()
	}
}

// BenchmarkListSelection_Toggle measures toggle selection performance
func BenchmarkListSelection_Toggle(b *testing.B) {
	items := convertToItems(createTestItems(1000))
	selMgr := list.NewSelectionManager(list.SelectionModeMultiple)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < len(items); i += 10 {
			selMgr.Toggle(items[i].ID())
		}
		selMgr.Clear()
	}
}

// BenchmarkViewportNavigation measures viewport navigation performance
func BenchmarkViewportNavigation(b *testing.B) {
	viewport := list.NewViewport(20, 1000)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < 100; i++ {
			viewport.MoveDown()
		}
		for i := 0; i < 100; i++ {
			viewport.MoveUp()
		}
		viewport.PageDown()
		viewport.PageUp()
		viewport.MoveToTop()
		viewport.MoveToBottom()
	}
}

// BenchmarkGroupManager measures group management performance
func BenchmarkGroupManager(b *testing.B) {
	groups := createTestGroups(100, 10)
	groupMgr := list.NewGroupManager()
	groupMgr.SetGroups(groups)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		groupMgr.ExpandAll()
		groupMgr.CollapseAll()
		groupMgr.ExpandAll()
	}
}

// Helper function to create test items
func createTestItems(count int) []list.FilterableItem {
	items := make([]list.FilterableItem, count)
	for i := 0; i < count; i++ {
		items[i] = list.NewListItem(
			fmt.Sprintf("id-%d", i),
			fmt.Sprintf("Item Title %d", i),
			fmt.Sprintf("Item Description %d", i),
		)
	}
	return items
}

// Helper function to create test groups
func createTestGroups(groupCount, itemsPerGroup int) []*list.Group {
	groups := make([]*list.Group, groupCount)
	for i := 0; i < groupCount; i++ {
		items := make([]list.Item, itemsPerGroup)
		for j := 0; j < itemsPerGroup; j++ {
			items[j] = list.NewListItem(
				fmt.Sprintf("group-%d-item-%d", i, j),
				fmt.Sprintf("Group %d Item %d", i, j),
				fmt.Sprintf("Description %d", j),
			)
		}
		groups[i] = list.NewGroup(
			fmt.Sprintf("Group %d", i),
			items,
		)
	}
	return groups
}

// Helper function to convert FilterableItem to Item
func convertToItems(filterable []list.FilterableItem) []list.Item {
	items := make([]list.Item, len(filterable))
	for i, item := range filterable {
		items[i] = item
	}
	return items
}
