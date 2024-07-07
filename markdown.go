package main

import (
	_ "embed"

	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/behavior"
	"github.com/richardwilkes/unison/example/demo"
)

// NewMarkdownWindow creates and displays our demo markdown window.
func NewMarkdownWindow(where unison.Point) (*unison.Window, error) {
	// Create the window
	wnd, err := unison.NewWindow("Gemiclip")
	if err != nil {
		return nil, err
	}

	// Install our menus
	installDefaultMenus(wnd)

	content := wnd.Content()
	content.SetLayout(&unison.FlexLayout{Columns: 1})
	content.SetBorder(
		unison.NewEmptyBorder(unison.NewSymmetricInsets(unison.StdHSpacing, unison.StdVSpacing)),
	)

	// Create the markdown view
	markdown := unison.NewMarkdown(true)
	markdown.SetContent(
		"",
		0,
	)
	go RunClipboardWatcher(markdown)

	// Create a scroll panel and place a table panel inside it
	scrollArea := unison.NewScrollPanel()
	scrollArea.SetContent(markdown, behavior.Fill, behavior.Fill)
	scrollArea.SetLayoutData(&unison.FlexLayoutData{
		HAlign: align.Fill,
		VAlign: align.Fill,
		HGrab:  true,
		VGrab:  true,
	})
	scrollArea.UpdateCursorCallback = func(_ unison.Point) *unison.Cursor {
		return unison.TextCursor()
	}

	content.AddChild(scrollArea)

	// Pack our window to fit its content, then set its location on the display and make it visible.
	wnd.Pack()
	rect := wnd.FrameRect()
	rect.Point = where
	wnd.SetFrameRect(rect)
	wnd.ToFront()
	wnd.ShowCursor()

	return wnd, nil
}

func installDefaultMenus(wnd *unison.Window) {
	unison.DefaultMenuFactory().BarForWindow(wnd, func(m unison.Menu) {
		unison.InsertStdMenus(m, demo.ShowAboutWindow, nil, nil)
		fileMenu := m.Menu(unison.FileMenuID)
		f := fileMenu.Factory()
		newMenu := f.NewMenu(demo.NewMenuID, "New…", nil)
		newMenu.InsertItem(-1, demo.NewWindowAction.NewMenuItem(f))
		newMenu.InsertItem(-1, demo.NewTableWindowAction.NewMenuItem(f))
		newMenu.InsertItem(-1, demo.NewDockWindowAction.NewMenuItem(f))
		newMenu.InsertItem(-1, demo.NewMarkdownWindowAction.NewMenuItem(f))
		fileMenu.InsertMenu(0, newMenu)
		fileMenu.InsertItem(1, demo.OpenAction.NewMenuItem(f))
		fileMenu.InsertItem(2, demo.ShowColorsWindowAction.NewMenuItem(f))
	})
}
