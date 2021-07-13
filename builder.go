/*
Package wgowut (wrapped gowut) provides convenient wrapper functions around the package "github.com/icza/gowut".
Initialize GUI components in one line by calling a make function and passing in a wgowut.Options
struct with any needed options. Unspecified options will be left at the default setting for the gwu
component. New options can/should be added but care should be taken to recognize the zero value of
the option type and the gwu default, so that options can be omitted and normal behavior occurs
and updates don't break existing GUIs (since defaults are respected). For examples,
see the MakeTable() CellPadding and HAlign as well as the MakeListBox() Enable option implementations.

Disclaimer

This documentation is not intended as a replacement for the gowut/gwu documentation; in order
to properly use wgowut, how to use gowut needs to be understood.

Recommended Usage

Create a struct in your application's GUI code that imports an anonymous *wgowut.GuiBuilder
struct. Your struct should also be used to store components needed for inputs etc. Prefer tables over
panels as it makes the code more readable and easy to understand. For the same reason, add high level
components to window or top level table/panel in order and at the same time. Example code:
 import (
	"github.com/ddrake12/wgowut"
	"github.com/icza/gowut/gwu"
 )

 type guiControl struct {
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
 	// make two more components


 	btnTable.Add(btn, 0, 0)
 	// add two more components in order to cells 0,1 and 0,2

	return btnTable
 }

*/
package wgowut

import (
	"strconv"

	"github.com/icza/gowut/gwu"
)

const (
	FullWidth  = "Full" // Use in Options.Width to set full width
	FullHeight = "Full" // Use in Options.Height to set full height
)

// Enable is used to set the Enable Option for gwu components that support it
type Enable int

const (
	EnableNil Enable = iota // Enable is used to set the Enable Option for gwu components that support it
	EnableTrue
	EnableFalse
)

// Layout is used to set the Layout Option for gwu components that support it
type Layout int

const (
	LayoutNil Layout = iota // Layout is used to set the Layout Option for gwu components that support it
	LayoutNatural
	LayoutHorizontal
	LayoutVertical
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
	// To actually see borders, BorderWidth and BorderStyle are required.
	BorderWidth              int
	BorderStyle, BorderColor string

	Layout            Layout // Layout is used for panels, tab panels, and tabbars and can be specified as Natural, Horizontal, or Vertical.
	Multi             bool
	Width, Height     string
	FontSize          string
	Color, Background string // Color is the 'foreground' color. For example, a label's text color is set using Color.
	ColSpan           int
	RowSpan           int
	Enable            Enable
	ReadOnly          bool
}

// NewGuiBuilder returns a GuiBuilder struct.
func NewGuiBuilder() *GuiBuilder {
	return &GuiBuilder{}
}

// MakeTable creates a gwu.Table and uses the following options:
//
// Rows, Cols, CellPadding, HAlign, Valign, Whitespace, BorderWidth, BorderStyle, BorderColor, Width, Height, FontSize, Color, Background
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

	setStyle(table.Style(), options)

	return table
}

func setStyle(style gwu.Style, options Options) {

	style.SetBorder2(options.BorderWidth, options.BorderStyle, options.BorderColor)

	if options.Width == FullWidth {
		style.SetFullWidth()
	} else if options.Width != "" {
		style.SetWidth(options.Width)
	}

	if options.Height == FullHeight {
		style.SetFullHeight()
	}
	if options.Height != "" {
		style.SetHeight(options.Height)
	}

	style.SetColor(options.Color)

	style.SetBackground(options.Background)

	style.SetWhiteSpace(options.WhiteSpace)

	style.SetFontSize(options.FontSize)
}

func setEnabled(comp gwu.HasEnabled, enable Enable) {
	switch enable {
	case EnableTrue:
		comp.SetEnabled(true)
	case EnableFalse:
		comp.SetEnabled(false)
	}
}

func setTableView(tView gwu.TableView, options Options) {
	tView.SetCellPadding(options.CellPadding)

	if options.HAlign != "" {
		tView.SetHAlign(options.HAlign)
	}
	if options.VAlign != "" {
		tView.SetVAlign(options.VAlign)
	}
}

func setLayout(pView gwu.PanelView, layout Layout) {
	switch layout {
	case LayoutNatural:
		pView.SetLayout(gwu.LayoutNatural)
	case LayoutHorizontal:
		pView.SetLayout(gwu.LayoutHorizontal)
	case LayoutVertical:
		pView.SetLayout(gwu.LayoutVertical)
	}
}

// FormatTableCell formats the given, table, row, and column and uses the following options:
//
// CellPadding, HAlign, VAlign, Whitespace, BorderWidth, BorderStyle, BorderColor, Width, Height, FontSize, Color, Background, ColSpan, RowSpan
func (g *GuiBuilder) FormatTableCell(table gwu.Table, row, col int, options Options) {

	padding := strconv.Itoa(options.CellPadding)
	table.CellFmt(row, col).Style().SetPadding(padding)

	if options.HAlign != "" {
		table.CellFmt(row, col).SetHAlign(options.HAlign)
	}
	if options.VAlign != "" {
		table.CellFmt(row, col).SetVAlign(options.VAlign)
	}

	table.SetColSpan(row, col, options.ColSpan)
	table.SetRowSpan(row, col, options.RowSpan)

	setStyle(table.CellFmt(row, col).Style(), options)

}

// MakeListBox takes in a slice of string values, adds them to a ListBox, and sets
// the first value to the default displayed/selected. The following options are
// used:
//
// Rows, Multi, BorderWidth, BorderStyle, BorderColor, Width, Height, FontSize, Color, Background, Enable
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

	setEnabled(lb, options.Enable)

	setStyle(lb.Style(), options)

	setStyle(lb.Style(), options)

	return lb
}

// MakeTextBox creates a text box with the given text.
// Note that the WhiteSpace option is only enforced if Enable is set to false or if ReadOnly is set to True.
// The following options are used:
//
// Rows, Cols, WhiteSpace BorderWidth, BorderStyle, BorderColor, Width, Height, FontSize, Color, Background, Enable, ReadOnly.
func (g *GuiBuilder) MakeTextBox(text string, options Options) gwu.TextBox {
	tb := gwu.NewTextBox(text)
	if options.Rows != 0 {
		tb.SetRows(options.Rows)
	}
	if options.Cols != 0 {
		tb.SetCols(options.Cols)
	}

	setEnabled(tb, options.Enable)

	tb.SetReadOnly(options.ReadOnly)

	setStyle(tb.Style(), options)

	return tb
}

// MakeLabel creates a label with the given text and uses following options:
//
// WhiteSpace, BorderWidth, BorderStyle, BorderColor, FontSize, Color, Background
func (g *GuiBuilder) MakeLabel(text string, options Options) gwu.Label {
	label := gwu.NewLabel(text)

	setStyle(label.Style(), options)

	return label
}

// MakeButton creates a button with the given text and uses the following options:
//
// WhiteSpace, BorderWidth, BorderStyle, BorderColor, Width, Height, FontSize, Color, Background
func (g *GuiBuilder) MakeButton(text string, options Options) gwu.Button {
	btn := gwu.NewButton(text)

	setStyle(btn.Style(), options)

	return btn
}

// MakeWindow creates a windows with the window list name and specific window/URL extension. Full width is always set.
// The following options are used:
//
// CellPadding, HAlign, VAlign, BorderWidth, BorderStyle, BorderColor, WhiteSpace, Color, Background
func (g *GuiBuilder) MakeWindow(name, extension string, options Options) gwu.Window {
	win := gwu.NewWindow(name, extension)

	setTableView(win, options)

	setStyle(win.Style(), options)

	return win
}

// MakePanel creates a gwu.Panel using the options.Layout parameter if specified. The following options are used:
//
// Layout, CellPadding, HAlign, Valign, WhiteSpace, BorderStyle, BorderWidth, BorderColor, Width, Height, Color, Background
func (g *GuiBuilder) MakePanel(options Options) gwu.Panel {

	panel := gwu.NewPanel()
	setLayout(panel, options.Layout)

	setTableView(panel, options)

	setStyle(panel.Style(), options)

	return panel
}

// AddLabelsToPanel creates a new gwu.Label with the given options for each labelText string, then adds them in order to a gwu.Panel.
func (g *GuiBuilder) AddLabelsToPanel(panel gwu.Panel, options Options, labelText ...string) {
	for _, text := range labelText {
		label := g.MakeLabel(text, options)
		panel.Add(label)
	}
}

// AddCompsToPanel adds a variable number of gwu.Comp interfaces to a gwu.Panel.
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

// MakeTabPanel creates a gwu.TabPanel using the options.Layout parameter if specified. The following options are used:
//
// Layout, CellPadding, HAlign, Valign, WhiteSpace, BorderStyle, BorderWidth, BorderColor, Width, Height, Color, Background
func (g *GuiBuilder) MakeTabPanel(options Options) gwu.TabPanel {

	tabPanel := gwu.NewTabPanel()

	setLayout(tabPanel, options.Layout)

	setTableView(tabPanel, options)

	setStyle(tabPanel.Style(), options)

	return tabPanel
}
