package page

type PageID string

// PageChangeMsg is used to change the current page
type PageChangeMsg struct {
	ID PageID
}

// PageCloseMsg is used to close the current page
type PageCloseMsg struct{}

// PageBackMsg is used to go back to the previous page
type PageBackMsg struct{}
