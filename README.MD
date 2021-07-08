# Introduction 
Wgowut (wrapped gowut) provides convenient wrapper functions around the package [gowut](github.com/icza/gowut). 

Initialize GUI components in one line by calling a make function and passing in a wgowut.Options struct with any needed options. 

Unspecified options will be left at the default setting for the gwu component. New options can/should be added but care should be taken to recognize the zero value of the option type and the gwu default, so that options can be omitted and normal behavior occurs and updates don't break existing GUIs (since defaults are respected). For examples, see the MakeTable() CellPadding and HAlign and the MakeListBox() Enable option implementations.

## Disclaimer
 
This documentation is not intended as a replacement for the gowut/gwu documentation; in order to properly use wgowut, how to use gowut needs to be understood.
 
## Recommended Usage

Create a struct in your application's GUI code that imports an anonymous *wgowut.GuiBuilder struct. Your struct should also be used to store components needed for inputs etc. Prefer tables over panels as it makes the code more readable and easy to understand. For the same reason, add high level components to window or top level table/panel in order and at the same time. Example code:    
###
```go
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
	return &guiControl{nil, nil, wgowut.NewGuiBuilder}
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
```

## Contributions Welcome! 