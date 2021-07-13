package wgowut

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/icza/gowut/gwu"
	"github.com/stretchr/testify/assert"
)

func checkTableView(t *testing.T, got gwu.TableView, options Options) {
	assert.Equal(t, options.CellPadding, got.CellPadding())
	assert.Equal(t, options.HAlign, got.HAlign())
	assert.Equal(t, options.VAlign, got.VAlign())
}

func checkStyle(t *testing.T, got gwu.Style, options Options) {

	if options.BorderWidth != 0 && options.BorderStyle != "" { //these are both required to actually display a border
		borderRe := regexp.MustCompile(`(\d+)\w\w\s+(\w+)\s+(\w+)`)
		matches := borderRe.FindStringSubmatch(got.Border())
		if matches != nil && len(matches) == 4 {
			assert.Equal(t, strconv.Itoa(options.BorderWidth), matches[1])
			assert.Equal(t, options.BorderStyle, matches[2])
			assert.Equal(t, options.BorderColor, matches[3])

		} else {
			t.Errorf(t.Name()+" checkStyle() - could not parse border information, revise test case. Border info: %v, regex: %v", got.Border(), borderRe.String())
		}
	}

	if options.Width == FullWidth {
		assert.Equal(t, "100%", got.Width())
	} else {
		assert.Equal(t, options.Width, got.Width())
	}

	assert.Equal(t, options.Height, got.Height())
	assert.Equal(t, options.Color, got.Color())
	assert.Equal(t, options.Background, got.Background())
	assert.Equal(t, options.WhiteSpace, got.WhiteSpace())
	assert.Equal(t, options.FontSize, got.FontSize())

}

func checkEnabled(t *testing.T, got gwu.HasEnabled, options Options) {
	if options.Enable == EnableTrue {
		assert.Equal(t, true, got.Enabled())
	} else if options.Enable == EnableFalse {
		assert.Equal(t, false, got.Enabled())
	}
}

func checkPanelView(t *testing.T, got gwu.PanelView, options Options) {
	switch options.Layout {
	case LayoutNatural:
		assert.Equal(t, gwu.LayoutNatural, got.Layout())
	case LayoutHorizontal:
		assert.Equal(t, gwu.LayoutHorizontal, got.Layout())
	case LayoutVertical:
		assert.Equal(t, gwu.LayoutVertical, got.Layout())
	}

}

func TestNewGuiBuilder(t *testing.T) {
	tests := []struct {
		name string
		want *GuiBuilder
	}{
		{"Constructor test", &GuiBuilder{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGuiBuilder()
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestGuiBuilder_MakeTable(t *testing.T) {

	tests := []struct {
		name    string
		options Options
	}{
		{"set all options", Options{
			Rows:        1,
			Cols:        1,
			CellPadding: 1,
			HAlign:      gwu.HARight,
			VAlign:      gwu.VABottom,
			WhiteSpace:  gwu.WhiteSpacePreWrap,
			BorderWidth: 2,
			BorderStyle: gwu.BrdStyleDotted,
			BorderColor: gwu.ClrFuchsia,
			Width:       "1",
			Height:      "1",
			FontSize:    "1",
			Color:       gwu.ClrMaroon,
			Background:  gwu.ClrAqua,
		}},
		{"set FullWidth and FullHeight", Options{Width: FullWidth, Height: FullHeight}},
		{"set no options", Options{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GuiBuilder{}
			got := g.MakeTable(tt.options)

			checkTableView(t, got.(gwu.TableView), tt.options)
			checkStyle(t, got.Style(), tt.options)
		})
	}
}

func TestGuiBuilder_FormatTableCell(t *testing.T) {

	tests := []struct {
		name    string
		options Options
	}{
		{"set all options", Options{
			CellPadding: 1,
			HAlign:      gwu.HARight,
			VAlign:      gwu.VABottom,
			WhiteSpace:  gwu.WhiteSpacePreWrap,
			BorderWidth: 2,
			BorderStyle: gwu.BrdStyleDotted,
			BorderColor: gwu.ClrFuchsia,
			Width:       "1",
			Height:      "1",
			FontSize:    "1",
			Color:       gwu.ClrMaroon,
			Background:  gwu.ClrAqua,
			ColSpan:     2,
			RowSpan:     2,
		}},
		{"set FullWidth and FullHeight", Options{Width: FullWidth, Height: FullHeight}},
		{"set no options", Options{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GuiBuilder{}
			table := g.MakeTable(Options{Rows: 5, Cols: 5})

			row, col := 1, 1
			g.FormatTableCell(table, row, col, tt.options)

			assert.Equal(t, strconv.Itoa(tt.options.CellPadding), table.CellFmt(row, col).Style().Padding())
			assert.Equal(t, tt.options.HAlign, table.CellFmt(row, col).HAlign())
			assert.Equal(t, tt.options.VAlign, table.CellFmt(row, col).VAlign())

			if tt.options.ColSpan != 0 {
				assert.Equal(t, tt.options.ColSpan, table.ColSpan(row, col))
			} else {
				assert.Equal(t, -1, table.ColSpan(row, col))
			}
			if tt.options.RowSpan != 0 {
				assert.Equal(t, tt.options.RowSpan, table.RowSpan(row, col))
			} else {
				assert.Equal(t, -1, table.RowSpan(row, col))
			}

			checkStyle(t, table.CellFmt(row, col).Style(), tt.options)

		})
	}
}

func TestGuiBuilder_MakeListBox(t *testing.T) {

	tests := []struct {
		name    string
		options Options
	}{
		{"set all options", Options{
			Rows:        3,
			WhiteSpace:  gwu.WhiteSpacePreWrap,
			BorderWidth: 2,
			BorderStyle: gwu.BrdStyleDotted,
			BorderColor: gwu.ClrFuchsia,
			Multi:       true,
			Width:       "1",
			Height:      "1",
			FontSize:    "1",
			Color:       gwu.ClrMaroon,
			Background:  gwu.ClrAqua,
			Enable:      EnableTrue,
		}},
		{"set FullWidth, FullHeight, and EnableFalse", Options{Width: FullWidth, Height: FullHeight, Enable: EnableFalse}},
		{"set no options", Options{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GuiBuilder{}
			values := []string{"hello", "world"}
			got := g.MakeListBox(values, tt.options)

			assert.Equal(t, values, got.Values())

			assert.Equal(t, tt.options.Rows, got.Rows())
			assert.Equal(t, tt.options.Multi, got.Multi())

			assert.Equal(t, true, got.Selected(0))

			checkEnabled(t, got.(gwu.HasEnabled), tt.options)
			checkStyle(t, got.Style(), tt.options)

		})
	}
}

func TestGuiBuilder_MakeTextBox(t *testing.T) {
	tests := []struct {
		name    string
		options Options
	}{
		{"set all options", Options{
			Rows:        3,
			Cols:        3,
			CellPadding: 1,
			WhiteSpace:  gwu.WhiteSpacePreWrap,
			BorderWidth: 2,
			BorderStyle: gwu.BrdStyleDotted,
			BorderColor: gwu.ClrFuchsia,
			Width:       "1",
			Height:      "1",
			FontSize:    "1",
			Color:       gwu.ClrMaroon,
			Background:  gwu.ClrAqua,
			Enable:      EnableTrue,
			ReadOnly:    true,
		}},
		{"set FullWidth, FullHeight, and EnableFalse", Options{Width: FullWidth, Height: FullHeight, Enable: EnableFalse}},
		{"set no options", Options{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GuiBuilder{}
			got := g.MakeTextBox(tt.name, tt.options)

			assert.Equal(t, tt.name, got.Text())

			if tt.options.Rows != 0 {
				assert.Equal(t, tt.options.Rows, got.Rows())
			}
			if tt.options.Cols != 0 {
				assert.Equal(t, tt.options.Cols, got.Cols())
			}

			checkEnabled(t, got.(gwu.HasEnabled), tt.options)

			assert.Equal(t, tt.options.ReadOnly, got.ReadOnly())

			checkStyle(t, got.Style(), tt.options)
		})
	}
}

func TestGuiBuilder_MakeLabel(t *testing.T) {

	tests := []struct {
		name    string
		options Options
	}{
		{"set all options", Options{
			WhiteSpace:  gwu.WhiteSpacePreWrap,
			BorderWidth: 2,
			BorderStyle: gwu.BrdStyleDotted,
			BorderColor: gwu.ClrFuchsia,
			Width:       "1",
			Height:      "1",
			FontSize:    "1",
			Color:       gwu.ClrMaroon,
			Background:  gwu.ClrAqua,
		}},
		{"set FullWidth and FullHeight", Options{Width: FullWidth, Height: FullHeight}},
		{"set no options", Options{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GuiBuilder{}
			got := g.MakeLabel(tt.name, tt.options)

			assert.Equal(t, tt.name, got.Text())

			checkStyle(t, got.Style(), tt.options)

		})
	}
}

func TestGuiBuilder_MakeButton(t *testing.T) {
	tests := []struct {
		name    string
		options Options
	}{
		{"set all options", Options{
			WhiteSpace:  gwu.WhiteSpacePreWrap,
			BorderWidth: 2,
			BorderStyle: gwu.BrdStyleDotted,
			BorderColor: gwu.ClrFuchsia,
			Width:       "1",
			Height:      "1",
			FontSize:    "1",
			Color:       gwu.ClrMaroon,
			Background:  gwu.ClrAqua,
		}},
		{"set FullWidth and FullHeight", Options{Width: FullWidth, Height: FullHeight}},
		{"set no options", Options{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GuiBuilder{}
			got := g.MakeButton(tt.name, tt.options)

			assert.Equal(t, tt.name, got.Text())
			checkStyle(t, got.Style(), tt.options)
		})
	}
}

func TestGuiBuilder_MakeWindow(t *testing.T) {
	tests := []struct {
		name    string
		options Options
	}{
		{"set all options", Options{
			CellPadding: 1,
			HAlign:      gwu.HARight,
			VAlign:      gwu.VABottom,
			WhiteSpace:  gwu.WhiteSpacePreWrap,
			BorderWidth: 2,
			BorderStyle: gwu.BrdStyleDotted,
			BorderColor: gwu.ClrFuchsia,
			Width:       "1",
			Height:      "1",
			FontSize:    "1",
			Color:       gwu.ClrMaroon,
			Background:  gwu.ClrAqua,
		}},
		{"set FullWidth and FullHeight", Options{Width: FullWidth, Height: FullHeight}},
		{"set no options", Options{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GuiBuilder{}
			got := g.MakeWindow(tt.name, "extension", tt.options)
			assert.Equal(t, tt.name, got.Name())

			checkTableView(t, got.(gwu.TableView), tt.options)

			checkStyle(t, got.Style(), tt.options)
		})
	}
}

func TestGuiBuilder_MakePanel(t *testing.T) {
	tests := []struct {
		name    string
		options Options
	}{
		{"set all options", Options{
			CellPadding: 1,
			HAlign:      gwu.HARight,
			VAlign:      gwu.VABottom,
			WhiteSpace:  gwu.WhiteSpacePreWrap,
			BorderWidth: 2,
			BorderStyle: gwu.BrdStyleDotted,
			BorderColor: gwu.ClrFuchsia,
			Layout:      LayoutHorizontal,
			Width:       "1",
			Height:      "1",
			FontSize:    "1",
			Color:       gwu.ClrMaroon,
			Background:  gwu.ClrAqua,
		}},
		{"set FullWidth, FullHeight, and LayoutNatural ", Options{Width: FullWidth, Height: FullHeight, Layout: LayoutNatural}},
		{"set LayoutVertical", Options{Layout: LayoutVertical}},
		{"set no options", Options{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GuiBuilder{}
			got := g.MakePanel(tt.options)

			checkTableView(t, got.(gwu.TableView), tt.options)

			checkPanelView(t, got.(gwu.PanelView), tt.options)

			checkStyle(t, got.Style(), tt.options)
		})
	}
}

func TestGuiBuilder_AddLabelsToPanel(t *testing.T) {

	tests := []struct {
		name    string
		options Options
		labels  []string
	}{
		{"set all options", Options{
			WhiteSpace:  gwu.WhiteSpacePreWrap,
			BorderWidth: 2,
			BorderStyle: gwu.BrdStyleDotted,
			BorderColor: gwu.ClrFuchsia,
			Width:       "1",
			Height:      "1",
			FontSize:    "1",
			Color:       gwu.ClrMaroon,
			Background:  gwu.ClrAqua,
		}, []string{"label 1, label 2"}},
		{"set FullWidth and FullHeight", Options{Width: FullWidth, Height: FullHeight}, []string{"label 1, label 2"}},
		{"set no options", Options{}, []string{"label 1, label 2"}},
		{"add no labels", Options{}, []string{""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GuiBuilder{}
			panel := g.MakePanel(Options{})

			g.AddLabelsToPanel(panel, tt.options, tt.labels...)
			for i, label := range tt.labels {
				got := panel.CompAt(i).(gwu.Label)
				assert.Equal(t, label, got.Text())
				checkStyle(t, got.Style(), tt.options)
			}
		})
	}
}

func TestGuiBuilder_AddCompsToPanel(t *testing.T) {

	tests := []struct {
		name  string
		comps []gwu.Comp
	}{
		{"add multiple comps", []gwu.Comp{gwu.NewButton("button"), gwu.NewLabel("label")}},
		{"add no comps", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GuiBuilder{}
			panel := g.MakePanel(Options{})
			g.AddCompsToPanel(panel, tt.comps...)

			if tt.comps == nil {
				got := panel.CompAt(0)
				assert.Equal(t, nil, got)
			} else {
				for i, comp := range tt.comps {
					got := panel.CompAt(i)
					assert.Equal(t, comp, got)
				}
			}
		})
	}
}

func TestGuiBuilder_SetEnabled(t *testing.T) {

	tests := []struct {
		name   string
		enable bool
		comps  []gwu.HasEnabled
	}{
		{"set multiple comps true", true, []gwu.HasEnabled{gwu.NewTextBox("text"), gwu.NewListBox([]string{"listbox"})}},
		{"set multiple comps false", false, []gwu.HasEnabled{gwu.NewTextBox("text"), gwu.NewListBox([]string{"listbox"})}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GuiBuilder{}
			g.SetEnabled(tt.enable, tt.comps...)
			for _, comp := range tt.comps {
				assert.Equal(t, tt.enable, comp.Enabled())
			}
		})
	}
}

func TestGuiBuilder_MakeTabPanel(t *testing.T) {
	tests := []struct {
		name    string
		options Options
	}{
		{"set all options", Options{
			CellPadding: 1,
			HAlign:      gwu.HARight,
			VAlign:      gwu.VABottom,
			WhiteSpace:  gwu.WhiteSpacePreWrap,
			BorderWidth: 2,
			BorderStyle: gwu.BrdStyleDotted,
			BorderColor: gwu.ClrFuchsia,
			Layout:      LayoutHorizontal,
			Width:       "1",
			Height:      "1",
			FontSize:    "1",
			Color:       gwu.ClrMaroon,
			Background:  gwu.ClrAqua,
		}},
		{"set FullWidth, FullHeight, and LayoutNatural ", Options{Width: FullWidth, Height: FullHeight, Layout: LayoutNatural}},
		{"set LayoutVertical", Options{Layout: LayoutVertical}},
		{"set no options", Options{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GuiBuilder{}
			got := g.MakeTabPanel(tt.options)

			checkTableView(t, got.(gwu.TableView), tt.options)

			checkPanelView(t, got.(gwu.PanelView), tt.options)
			checkStyle(t, got.Style(), tt.options)
		})
	}
}
