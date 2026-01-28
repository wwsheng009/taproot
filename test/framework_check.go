package main

import (
	"fmt"
	"os"

	"github.com/yourorg/taproot/internal/tui/app"
	"github.com/yourorg/taproot/internal/tui/components/dialogs"
	"github.com/yourorg/taproot/internal/tui/page"
	"github.com/yourorg/taproot/internal/tui/util"
	"github.com/yourorg/taproot/internal/tui/components/logo"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	pageHome page.PageID = "home"
	pageTest page.PageID = "test"
)

func main() {
	// Test 1: Create application
	application := app.NewApp()
	fmt.Println("✓ AppModel created")

	// Test 2: Register pages
	homePage := HomePage{}
	testPage := TestPage{}
	application.RegisterPage(pageHome, homePage)
	application.RegisterPage(pageTest, testPage)
	fmt.Println("✓ Pages registered")

	// Test 3: Set initial page
	application.SetPage(pageHome)
	fmt.Println("✓ Initial page set")

	// Test 4: Test logo rendering
	logoText := logo.SmallRender(80)
	if logoText == "" {
		fmt.Println("✗ Logo rendering failed")
		os.Exit(1)
	}
	fmt.Println("✓ Logo rendering works")

	// Test 5: Test dialog creation
	dialog := &TestDialog{}
	if dialog.ID() == "" {
		fmt.Println("✗ Dialog ID is empty")
		os.Exit(1)
	}
	fmt.Println("✓ Dialog creation works")

	// Test 6: Test page stack navigation
	msg := page.PageChangeMsg{ID: pageTest}
	newModel, _ := application.Update(msg)
	application = newModel.(app.AppModel)
	if application.CurrentPage() != pageTest {
		fmt.Println("✗ Page navigation failed, got:", application.CurrentPage())
		os.Exit(1)
	}
	fmt.Println("✓ Page navigation works")

	// Test 7: Test page back
	backMsg := page.PageBackMsg{}
	newModel, _ = application.Update(backMsg)
	application = newModel.(app.AppModel)
	if application.CurrentPage() != pageHome {
		fmt.Println("✗ Page back failed, got:", application.CurrentPage())
		os.Exit(1)
	}
	fmt.Println("✓ Page back navigation works")

	// Test 8: Test dialog open
	openDialogMsg := dialogs.OpenDialogMsg{Model: dialog}
	newModel, _ = application.Update(openDialogMsg)
	application = newModel.(app.AppModel)
	if !application.HasDialogs() {
		fmt.Println("✗ Dialog open failed")
		os.Exit(1)
	}
	fmt.Println("✓ Dialog open works")

	// Test 9: Test dialog close
	closeDialogMsg := dialogs.CloseDialogMsg{}
	newModel, _ = application.Update(closeDialogMsg)
	application = newModel.(app.AppModel)
	if application.HasDialogs() {
		fmt.Println("✗ Dialog close failed")
		os.Exit(1)
	}
	fmt.Println("✓ Dialog close works")

	fmt.Println("\n🎉 All framework tests passed!")
	fmt.Println("\nTo run the interactive demo:")
	fmt.Println("  go run examples/app/main.go")
}

// HomePage is a simple home page
type HomePage struct{}

func (h HomePage) Init() tea.Cmd { return nil }
func (h HomePage) Update(msg tea.Msg) (util.Model, tea.Cmd) { return h, nil }
func (h HomePage) View() string {
	t := lipgloss.NewStyle()
	return t.Bold(true).Foreground(lipgloss.Color("86")).Render("Home Page")
}

// TestPage is a simple test page
type TestPage struct{}

func (t TestPage) Init() tea.Cmd { return nil }
func (t TestPage) Update(msg tea.Msg) (util.Model, tea.Cmd) { return t, nil }
func (t TestPage) View() string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("Test Page")
}

// TestDialog is a simple test dialog
type TestDialog struct{}

func (d *TestDialog) Init() tea.Cmd { return nil }
func (d *TestDialog) Update(msg tea.Msg) (util.Model, tea.Cmd) { return d, nil }
func (d *TestDialog) View() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("white")).
		Background(lipgloss.Color("62")).
		Padding(1, 2).
		Render("Test Dialog")
}
func (d *TestDialog) Position() (int, int) { return 10, 5 }
func (d *TestDialog) ID() dialogs.DialogID { return "test-dialog" }
