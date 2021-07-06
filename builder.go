/*
Package wgowut (wrapped gowut) provides convenient wrapper functions around the package "github.com/icza/gowut".
Initialize GUI components in one line by calling a make function and passing in a wgowut.Options
struct with any needed options. Unspecified options will be left at the default setting for the gwu
component. New options can/should be added but care should be taken to recognize the zero value of
the option type and the gwu default, so that options can be omitted and normal behavior occurs
and updates don't break existing GUIs (since defaults are respected). For examples,
see the MakeTable() CellPadding and HAlign option implementations.

Disclaimer

This documentation is not intended as a replacement for the gowut/gwu documentation; in order
to properly use wgowut, how to use gowut needs to be understood.

Recommended Usage

Create a struct in your application's GUI code that imports an anonymous *wgowut.GuiBuilder
struct. Your struct should also be used to store components needed for inputs etc. Prefer tables over
panels as it makes the code more readable and easy to understand. For the same reason, add high level
components to window or top level table/panel in order and at the same time. Example code:
 import "github.com/ddrake12/wgowut"

 struct guiControl {
	importantTb gwu.TextBox
	importantLb gwu.ListBox
	*wgowut.GuiBuilder
 }

 func newGuiControl() *guiControl {
	return &guiControl{nil, nil, wgowut.NewGuiBuilder()}
 }

 func StartGui() {
	gc := newGuiControl()

	win := gc.MakeWindow("urlExtension", "application", wgowut.Options{CellPadding: 10})
	btnTable := gc.makeBtnTable()
	inputTable := gc.makeInputTable() // Not shown, but here guiControl.importantTb and guiControl.importantLb would be created
	// make more stuff

	// add components to window or top level table/panel in order:
	win.Add(inputTable)
	win.Add(btnTable) // btnTable on bottom if last added component to a gwu.Window

	// start gwu server
 }

 func (gc *guiControl) makeBtnTable() gwu.Table {

 	btnTable := gc.MakeTable(wgowut.Options{Rows: 1, Cols: 3, CellPadding: 5, HAlign: gwu.HARight})
 	btn := gwu.NewButton("Start")
 	btn.AddEHandlerFunc(func(e gwu.Event) {
		currentText := gc.importantTb.Text()
		selectedVal := gc.importantLb.SelectedValue()
		// do something with these values
 	}, gwu.ETypeClick)
 	// make two more componenets

 	btnTable.Add(btn, 0, 0)
 	// add two more components in order to cells 0,1 and 0,2

	return btnTable
 }

*/
package wgowut

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/icza/gowut/gwu"
)

// GuiBuilder is an empty struct that allows convenient access to package functions.
type GuiBuilder struct {
}

// Options implements flags for standard gwu options used while creating components. These options are not required and the
// gwu default will be used when the option is left blank. Some gwu types inherit certain attributes from a parent's style, some don't.
// Details about certain options are shown below.
type Options struct {
	Rows, Cols  int
	CellPadding int
	HAlign      gwu.HAlign
	VAlign      gwu.VAlign
	WhiteSpace  string
	// To make borders, BorderWidth and BorderStyle are required.
	BorderWidth              int
	BorderStyle, BorderColor string

	Orientation       string // Orientation is used for panels and can be specified as Natural or Horizontal. Default orientation is vertical.
	Multi             bool
	Width, Height     string
	FontSize          string
	Color, Background string // Color is the 'foreground' color. For example, a label's text color is set using Color.
	ColSpan           int
	RowSpan           int
	Enable            string
	ReadOnly          bool
}

// NewGuiBuilder returns a GuiBuilder struct.
func NewGuiBuilder() *GuiBuilder {
	return &GuiBuilder{}
}

// MakeTable creates a gwu.Table and accepts the following options:
//
// Rows, Cols, CellPadding, HAlign, Valign, BorderWidth, BorderStyle, BorderColor, Width, Height, Color, Background
func (g *GuiBuilder) MakeTable(options Options) gwu.Table {
	table := gwu.NewTable()

	table.EnsureSize(options.Rows, options.Cols)

	table.SetCellPadding(options.CellPadding)

	if options.HAlign != "" {
		table.SetHAlign(options.HAlign)
	}
	if options.VAlign != "" {
		table.SetVAlign(options.VAlign)
	}

	table.Style().SetWhiteSpace(options.WhiteSpace)

	if options.BorderWidth != 0 || options.BorderStyle != "" || options.BorderColor != "" {
		if options.BorderWidth != 0 && options.BorderStyle != "" {
			table.Style().SetBorder2(options.BorderWidth, options.BorderStyle, options.BorderColor)
		} else {
			fmt.Fprintf(os.Stderr, "\nError: Setting a border requires style and width.\n")
		}
	}

	if strings.EqualFold(options.Width, "Full") && strings.EqualFold(options.Height, "Full") {
		table.Style().SetFullWidth()
		table.Style().SetFullHeight()
	} else if strings.EqualFold(options.Width, "Full") {
		table.Style().SetFullWidth()
		table.Style().SetHeight(options.Height)
	} else if strings.EqualFold(options.Height, "Full") {
		table.Style().SetFullHeight()
		table.Style().SetWidth(options.Width)
	} else {
		table.Style().SetSize(options.Width, options.Height)
	}

	table.Style().SetColor(options.Color)

	table.Style().SetBackground(options.Background)

	return table
}

// FormatTableCell formats the given, table, row, and column. The following options are accepted:
//
// CellPadding, HAlign, VAlign, WhiteSpace, BorderWidth, BorderStyle, BorderColor, Width, Height, Color, Background, ColSpan
func (g *GuiBuilder) FormatTableCell(table gwu.Table, row, col int, options Options) {

	padding := strconv.Itoa(options.CellPadding)
	table.CellFmt(row, col).Style().SetPadding(padding)

	if options.HAlign != "" {
		table.CellFmt(row, col).SetHAlign(options.HAlign)
	}
	if options.VAlign != "" {
		table.CellFmt(row, col).SetVAlign(options.VAlign)
	}

	if options.BorderWidth != 0 || options.BorderStyle != "" || options.BorderColor != "" {
		if options.BorderWidth != 0 && options.BorderStyle != "" {
			table.CellFmt(row, col).Style().SetBorder2(options.BorderWidth, options.BorderStyle, options.BorderColor)
		} else {
			fmt.Fprintf(os.Stderr, "\nError: Setting a border requires style and width.\n")
		}
	}

	if strings.EqualFold(options.Width, "Full") && strings.EqualFold(options.Height, "Full") {
		table.CellFmt(row, col).Style().SetFullWidth()
		table.CellFmt(row, col).Style().SetFullHeight()
	} else if strings.EqualFold(options.Width, "Full") {
		table.CellFmt(row, col).Style().SetFullWidth()
		table.CellFmt(row, col).Style().SetHeight(options.Height)
	} else if strings.EqualFold(options.Height, "Full") {
		table.CellFmt(row, col).Style().SetFullHeight()
		table.CellFmt(row, col).Style().SetWidth(options.Width)
	} else {
		table.CellFmt(row, col).Style().SetSize(options.Width, options.Height)
	}

	table.CellFmt(row, col).Style().SetColor(options.Color)

	table.CellFmt(row, col).Style().SetBackground(options.Background)

	table.SetColSpan(row, col, options.ColSpan)

	table.SetRowSpan(row, col, options.RowSpan)

}

// MakeListBox takes in a slice of string values, adds them to a ListBox, and sets
// the first value to the default displayed/selected. The following options are
// accepted:
//
// Rows, Multi, Width, Height, FontSize, Color, Background, Enable
func (g *GuiBuilder) MakeListBox(values []string, options Options) gwu.ListBox {
	lb := gwu.NewListBox(values)

	lb.SetRows(options.Rows) // technically this zero value doesn't match the gwu default, but the
	// default is unusable IMO with multi and this decision was made pre-wgowut release

	if options.Multi {
		lb.SetMulti(true)
	}
	if len(values) != 0 {
		lb.SetSelected(0, true)
	}

	if strings.EqualFold(options.Width, "Full") && strings.EqualFold(options.Height, "Full") {
		lb.Style().SetFullWidth()
		lb.Style().SetFullHeight()
	} else if strings.EqualFold(options.Width, "Full") {
		lb.Style().SetFullWidth()
		lb.Style().SetHeight(options.Height)
	} else if strings.EqualFold(options.Height, "Full") {
		lb.Style().SetFullHeight()
		lb.Style().SetWidth(options.Width)
	} else {
		lb.Style().SetSize(options.Width, options.Height)
	}

	lb.Style().SetFontSize(options.FontSize)

	lb.Style().SetColor(options.Color)

	lb.Style().SetBackground(options.Background)

	if strings.EqualFold(options.Enable, "true") {
		lb.SetEnabled(true)
	}
	if strings.EqualFold(options.Enable, "false") {
		lb.SetEnabled(false)
	}

	return lb
}

// MakeTextBox creates a text box with the given text. The following options are accepted:
//
// Rows, Cols, BorderWidth, BorderStyle, BorderColor, Width, Height, FontSize, Color, Background, Enable, ReadOnly
func (g *GuiBuilder) MakeTextBox(text string, options Options) gwu.TextBox {
	tb := gwu.NewTextBox(text)
	if options.Rows != 0 {
		tb.SetRows(options.Rows)
	}
	if options.Cols != 0 {
		tb.SetCols(options.Cols)
	}

	if options.BorderWidth != 0 || options.BorderStyle != "" || options.BorderColor != "" {
		if options.BorderWidth != 0 && options.BorderStyle != "" {
			tb.Style().SetBorder2(options.BorderWidth, options.BorderStyle, options.BorderColor)
		} else {
			fmt.Fprintf(os.Stderr, "\nError: Setting a border requires style and width.\n")
		}
	}

	if strings.EqualFold(options.Width, "Full") && strings.EqualFold(options.Height, "Full") {
		tb.Style().SetFullWidth()
		tb.Style().SetFullHeight()
	} else if strings.EqualFold(options.Width, "Full") {
		tb.Style().SetFullWidth()
		tb.Style().SetHeight(options.Height)
	} else if strings.EqualFold(options.Height, "Full") {
		tb.Style().SetFullHeight()
		tb.Style().SetWidth(options.Width)
	} else {
		tb.Style().SetSize(options.Width, options.Height)
	}

	tb.Style().SetFontSize(options.FontSize)

	tb.Style().SetColor(options.Color)

	tb.Style().SetBackground(options.Background)

	if strings.EqualFold(options.Enable, "true") {
		tb.SetEnabled(true)
	}
	if strings.EqualFold(options.Enable, "false") {
		tb.SetEnabled(false)
	}

	tb.SetReadOnly(options.ReadOnly)

	return tb
}

// MakeLabel creates a label with the given text. The following options are accepted:
//
// WhiteSpace, BorderWidth, BorderStyle, BorderColor, FontSize, Color, Background
func (g *GuiBuilder) MakeLabel(text string, options Options) gwu.Label {
	label := gwu.NewLabel(text)

	label.Style().SetWhiteSpace(options.WhiteSpace)

	if options.BorderWidth != 0 || options.BorderStyle != "" || options.BorderColor != "" {
		if options.BorderWidth != 0 && options.BorderStyle != "" {
			label.Style().SetBorder2(options.BorderWidth, options.BorderStyle, options.BorderColor)
		} else {
			fmt.Fprint(os.Stderr, "\nError: Setting a border requires style and width.\n")
		}
	}

	label.Style().SetFontSize(options.FontSize)

	//Height and Width do nothing on labels, cannot implement

	label.Style().SetColor(options.Color)

	label.Style().SetBackground(options.Background)

	return label
}

// MakeButton creates a button with the given text. The following options are accepted:
//
// WhiteSpace, BorderWidth, BorderStyle, BorderColor, Width, Height, FontSize, Color, Background
func (g *GuiBuilder) MakeButton(text string, options Options) gwu.Button {
	btn := gwu.NewButton(text)

	btn.Style().SetWhiteSpace(options.WhiteSpace)

	if options.BorderWidth != 0 || options.BorderStyle != "" || options.BorderColor != "" {
		if options.BorderWidth != 0 && options.BorderStyle != "" {
			btn.Style().SetBorder2(options.BorderWidth, options.BorderStyle, options.BorderColor)
		} else {
			fmt.Fprintf(os.Stderr, "\nError: Setting a border requires style and width.\n")
		}
	}

	if strings.EqualFold(options.Width, "Full") && strings.EqualFold(options.Height, "Full") {
		btn.Style().SetFullWidth()
		btn.Style().SetFullHeight()
	} else if strings.EqualFold(options.Width, "Full") {
		btn.Style().SetFullWidth()
		btn.Style().SetHeight(options.Height)
	} else if strings.EqualFold(options.Height, "Full") {
		btn.Style().SetFullHeight()
		btn.Style().SetWidth(options.Width)
	} else {
		btn.Style().SetSize(options.Width, options.Height)
	}

	btn.Style().SetFontSize(options.FontSize)

	btn.Style().SetColor(options.Color)

	btn.Style().SetBackground(options.Background)

	return btn
}

// MakeWindow creates a windows with the window list name and specific window/URL extension. Full width is always set.
// The following options are accepted:
//
// CellPadding, HAlign, VAlign, BorderWidth, BorderStyle, BorderColor, WhiteSpace, Color, Background
func (g *GuiBuilder) MakeWindow(name, extension string, options Options) gwu.Window {
	win := gwu.NewWindow(name, extension)

	win.SetCellPadding(options.CellPadding)

	if options.HAlign != "" {
		win.SetHAlign(options.HAlign)
	}
	if options.VAlign != "" {
		win.SetVAlign(options.VAlign)
	}

	if options.Width != "" || options.Height != "" {
		win.Style().SetSize(options.Width, options.Height)
	} else {
		win.Style().SetFullWidth()
	}

	win.Style().SetWhiteSpace(options.WhiteSpace)

	if options.BorderWidth != 0 || options.BorderStyle != "" || options.BorderColor != "" {
		if options.BorderWidth != 0 && options.BorderStyle != "" {
			win.Style().SetBorder2(options.BorderWidth, options.BorderStyle, options.BorderColor)
		} else {
			fmt.Fprintf(os.Stderr, "\nError: Setting a border requires style and width.\n")
		}
	}

	win.Style().SetColor(options.Color)

	win.Style().SetBackground(options.Background)

	return win
}

// MakePanel creates a gwu.Panel that is default veritcal orientation. The orientation paramter can also be specified as Natural or Horizontal. The following options are accepted:
//
// Orientation, CellPadding, HAlign, Valign, WhiteSpace, BorderStyle, BorderWidth, BorderColor, Width, Height, Color, Background
func (g *GuiBuilder) MakePanel(options Options) gwu.Panel {
	var panel gwu.Panel

	if strings.EqualFold(options.Orientation, "Natural") {
		panel = gwu.NewNaturalPanel()
	} else if strings.EqualFold(options.Orientation, "Horizontal") {
		panel = gwu.NewHorizontalPanel()
	} else {
		panel = gwu.NewVerticalPanel()
	}

	panel.SetCellPadding(options.CellPadding)

	if options.HAlign != "" {
		panel.SetHAlign(options.HAlign)
	}
	if options.VAlign != "" {
		panel.SetVAlign(options.VAlign)
	}

	panel.Style().SetWhiteSpace(options.WhiteSpace)

	if options.BorderWidth != 0 || options.BorderStyle != "" || options.BorderColor != "" {
		if options.BorderWidth != 0 && options.BorderStyle != "" {
			panel.Style().SetBorder2(options.BorderWidth, options.BorderStyle, options.BorderColor)
		} else {
			fmt.Fprintf(os.Stderr, "\nError: Setting a border requires style and width.\n")
		}
	}

	if strings.EqualFold(options.Width, "Full") && strings.EqualFold(options.Height, "Full") {
		panel.Style().SetFullWidth()
		panel.Style().SetFullHeight()
	} else if strings.EqualFold(options.Width, "Full") {
		panel.Style().SetFullWidth()
		panel.Style().SetHeight(options.Height)
	} else if strings.EqualFold(options.Height, "Full") {
		panel.Style().SetFullHeight()
		panel.Style().SetWidth(options.Width)
	} else {
		panel.Style().SetSize(options.Width, options.Height)
	}

	panel.Style().SetColor(options.Color)

	panel.Style().SetBackground(options.Background)

	return panel
}

// AddLabelsToPanel creates a new gwu.Label for each string, adding them in order to a gwu.Panel.
func (g *GuiBuilder) AddLabelsToPanel(panel gwu.Panel, labels ...string) {
	for _, label := range labels {
		panel.Add(gwu.NewLabel(label))
	}
}

// AddsCompsToPanel adds a variable number of gwu.Comp interfaces to a gwu.Panel.
func (g *GuiBuilder) AddCompsToPanel(panel gwu.Panel, comps ...gwu.Comp) {
	for _, comp := range comps {
		panel.Add(comp)
	}
}

// SetEnabled sets enabled on a variable number of gwu.HasEnabled interfaces
func (g *GuiBuilder) SetEnabled(enable bool, comps ...gwu.HasEnabled) {
	for _, comp := range comps {
		comp.SetEnabled(enable)
	}
}

// MakeTabPanel creates a gwu.TanPanel that is default veritcal orientation. The orientation paramter can also be specified as Natural or Horizontal. The following options are accepted:
//
// Orientation, CellPadding, HAlign, Valign, BorderStyle, BorderWidth, BorderColor, Width, Height, Color, Background
func (g *GuiBuilder) MakeTabPanel(options Options) gwu.TabPanel {

	tabPanel := gwu.NewTabPanel()

	tabPanel.SetCellPadding(options.CellPadding)

	if options.HAlign != "" {
		tabPanel.SetHAlign(options.HAlign)
	}
	if options.VAlign != "" {
		tabPanel.SetVAlign(options.VAlign)
	}

	if options.BorderWidth != 0 || options.BorderStyle != "" || options.BorderColor != "" {
		if options.BorderWidth != 0 && options.BorderStyle != "" {
			tabPanel.Style().SetBorder2(options.BorderWidth, options.BorderStyle, options.BorderColor)
		} else {
			fmt.Fprintf(os.Stderr, "\nError: Setting a border requires style and width.\n")
		}
	}

	if strings.EqualFold(options.Width, "Full") && strings.EqualFold(options.Height, "Full") {
		tabPanel.Style().SetFullWidth()
		tabPanel.Style().SetFullHeight()
	} else if strings.EqualFold(options.Width, "Full") {
		tabPanel.Style().SetFullWidth()
		tabPanel.Style().SetHeight(options.Height)
	} else if strings.EqualFold(options.Height, "Full") {
		tabPanel.Style().SetFullHeight()
		tabPanel.Style().SetWidth(options.Width)
	} else {
		tabPanel.Style().SetSize(options.Width, options.Height)
	}

	tabPanel.Style().SetColor(options.Color)

	tabPanel.Style().SetBackground(options.Background)

	return tabPanel
}
