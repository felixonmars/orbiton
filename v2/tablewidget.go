package main

import (
	"fmt"
	"strings"

	"github.com/xyproto/vt100"
)

// TableWidget represents a TUI widget for editing a Markdown table
type TableWidget struct {
	title          string               // title
	contents       *[][]string          // the table contents
	bgColor        vt100.AttributeColor // background color
	highlightColor vt100.AttributeColor // selected color (the choice that has been selected after return has been pressed)
	headerColor    vt100.AttributeColor // the color of the table header row
	textColor      vt100.AttributeColor // text color (the choices that are not highlighted)
	titleColor     vt100.AttributeColor // title color (above the choices)
	cursorColor    vt100.AttributeColor // color of the "_" cursor
	commentColor   vt100.AttributeColor // comment color
	cx             int                  // current content position
	marginLeft     int                  // margin, may be negative?
	marginTop      int                  // margin, may be negative?
	oldy           int                  // previous position
	cy             int                  // current content position
	oldx           int                  // previous position
	h              int                  // height (number of menu items)
	w              int                  // width
}

// NewTableWidget creates a new TableWidget
func NewTableWidget(title string, contents *[][]string, titleColor, headerColor, textColor, highlightColor, cursorColor, commentColor, bgColor vt100.AttributeColor, canvasWidth, canvasHeight, initialY int) *TableWidget {

	columnWidths := TableColumnWidths([]string{}, *contents)

	widgetWidth := 0
	for _, w := range columnWidths {
		widgetWidth += w + 1
	}
	if widgetWidth > int(canvasWidth) {
		widgetWidth = int(canvasWidth)
	}

	widgetHeight := len(*contents)

	return &TableWidget{
		title:          title,
		w:              widgetWidth,
		h:              widgetHeight,
		cx:             0,
		oldx:           0,
		cy:             initialY,
		oldy:           initialY,
		marginLeft:     10,
		marginTop:      10,
		contents:       contents,
		titleColor:     titleColor,
		headerColor:    headerColor,
		textColor:      textColor,
		highlightColor: highlightColor,
		cursorColor:    cursorColor,
		commentColor:   commentColor,
		bgColor:        bgColor,
	}
}

// Expand the table contents to the longest width
func Expand(contents *[][]string) {
	// Find the max width
	maxWidth := 0
	for y := 0; y < len(*contents); y++ {
		if len((*contents)[y]) > maxWidth {
			maxWidth = len(*contents)
		}
	}
	// Find all rows less than max width
	for y := 0; y < len(*contents); y++ {
		if (*contents)[y] == nil {
			// Initialize the row
			(*contents)[y] = make([]string, maxWidth)
		} else if len((*contents)[y]) < maxWidth {
			backup := (*contents)[y]
			// Expand the row by creating a blank string slice
			(*contents)[y] = make([]string, maxWidth)
			// Fill in the old data for the first fields of the row
			copy((*contents)[y], backup)
		}
	}
}

// ContentsWH returns the width and the height of the table contents
func (tw *TableWidget) ContentsWH() (int, int) {
	rowCount := len(*tw.contents)
	if rowCount == 0 {
		return 0, 0
	}
	return len((*tw.contents)[0]), rowCount
}

// Draw will draw this menu widget on the given canvas
func (tw *TableWidget) Draw(c *vt100.Canvas) {
	cw, ch := tw.ContentsWH()

	// Height of the title + the size + a blank line
	titleHeight := 3

	// Draw the title
	title := tw.title
	for x, r := range title {
		c.PlotColor(uint(tw.marginLeft+x), uint(tw.marginTop), tw.titleColor, r)
	}

	// Plot the table size below the title
	sizeString := fmt.Sprintf("%dx%d", cw, ch)
	for x, r := range sizeString {
		c.PlotColor(uint(tw.marginLeft+x), uint(tw.marginTop+1), tw.commentColor, r)
	}

	columnWidths := TableColumnWidths([]string{}, *tw.contents)

	// Draw the headers, with various colors
	// Draw the menu entries, with various colors
	for y := 0; y < ch; y++ {
		xpos := tw.marginLeft
		// First clear this row with spaces
		spaces := strings.Repeat(" ", int(c.W()))
		c.Write(0, uint(tw.marginTop+y+titleHeight), tw.textColor, tw.bgColor, spaces)
		for x := 0; x < len((*tw.contents)[y]); x++ {
			field := (*tw.contents)[y][x]
			color := tw.textColor
			if y == int(tw.cy) && x == int(tw.cx) {
				color = tw.highlightColor
				// Draw the "cursor"
				c.Write(uint(xpos+len(field)), uint(tw.marginTop+y+titleHeight), tw.cursorColor, tw.bgColor, "_")
			} else if y == 0 {
				color = tw.headerColor
			}
			c.Write(uint(xpos), uint(tw.marginTop+y+titleHeight), color, tw.bgColor, field)
			xpos += columnWidths[x] + 2
		}
	}

	// Clear four extra rows after the table
	spaces := strings.Repeat(" ", int(c.W()))
	for y := ch; y < ch+4; y++ {
		if uint(y) < c.H() {
			c.Write(0, uint(tw.marginTop+y+titleHeight), tw.textColor, tw.bgColor, spaces)
		}
	}

	// Plot the table size below the table, centered
	//sizeString := fmt.Sprintf("%dx%d", cw, ch)
	//for x, r := range sizeString {
	//c.PlotColor((c.W()/2)-uint(len(sizeString))+uint(x), uint(tw.marginTop+ch+titleHeight+2), tw.cursorColor, r)
	//}
}

// Up will move the highlight up (with wrap-around)
func (tw *TableWidget) Up() {
	cw, ch := tw.ContentsWH()

	tw.oldy = tw.cy
	tw.cy--
	if tw.cy < 0 {
		tw.cy = ch - 1
	}
	// just in case rows have differing lengths
	if tw.cx >= cw {
		tw.cx = cw - 1
	}
}

// Down will move the highlight down (with wrap-around)
func (tw *TableWidget) Down() {
	cw, ch := tw.ContentsWH()

	tw.oldy = tw.cy
	tw.cy++
	if tw.cy >= ch {
		tw.cy = 0
	}
	// just in case rows have differing lengths
	if tw.cx >= cw {
		tw.cx = cw - 1
	}
}

// Left will move the highlight left (with wrap-around)
func (tw *TableWidget) Left() {
	cw, _ := tw.ContentsWH()

	tw.oldx = tw.cx
	tw.cx--
	if tw.cx < 0 {
		tw.cx = cw - 1
	}
}

// Right will move the highlight right (with wrap-around)
func (tw *TableWidget) Right() {
	cw, _ := tw.ContentsWH()

	tw.oldx = tw.cx
	tw.cx++
	if tw.cx >= cw {
		tw.cx = 0
	}
}

// NextOrInsert will move the highlight to the next cell, or insert a new row
func (tw *TableWidget) NextOrInsert() {
	cw, ch := tw.ContentsWH()
	tw.oldx = tw.cx
	tw.cx++
	if tw.cx >= cw {
		tw.cx = 0
		tw.cy++
		if tw.cy >= ch {
			newRow := make([]string, cw)
			(*tw.contents) = append((*tw.contents), newRow)
			tw.h++     // Update the widget table height as well (this is not the content height)
			tw.cy = ch // old max index + 1
		}
	}
}

// InsertRowBelow will insert a row below this one
func (tw *TableWidget) InsertRowBelow() {
	cw, _ := tw.ContentsWH()
	tw.cx = 0
	tw.cy++

	newRow := make([]string, cw)
	// Insert the new row at the cy position
	*tw.contents = append((*tw.contents)[:tw.cy], append([][]string{newRow}, (*tw.contents)[tw.cy:]...)...)

	tw.h++ // Update the widget table height as well (this is not the content height)
}

// CurrentRowIsEmpty checks if the current row is empty
func (tw *TableWidget) CurrentRowIsEmpty() bool {
	row := (*tw.contents)[tw.cy]
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

// DeleteCurrentRow deletes the current row
func (tw *TableWidget) DeleteCurrentRow() {
	if tw.cy >= 0 && tw.cy < len(*tw.contents) {
		// Remove the current row from the contents
		*tw.contents = append((*tw.contents)[:tw.cy], (*tw.contents)[tw.cy+1:]...)
		tw.cy--
		tw.h-- // Update the widget table height as well (this is not the content height)
	}
}

// SelectIndex will select a specific index. Returns false if it was not possible.
func (tw *TableWidget) SelectIndex(x, y int) bool {
	cw, ch := tw.ContentsWH()

	if x >= cw || y >= ch {
		return false
	}
	tw.oldx = tw.cx
	tw.oldy = tw.cy
	tw.cx = x
	tw.cy = y
	return true
}

// SelectStart will select the start of the row
func (tw *TableWidget) SelectStart() bool {
	return tw.SelectIndex(0, tw.cy)
}

// SelectEnd will select the start of the row
func (tw *TableWidget) SelectEnd() bool {
	cw, _ := tw.ContentsWH()
	return tw.SelectIndex(cw-1, tw.cy)
}

// Set will change the field contents of the current position
func (tw *TableWidget) Set(field string) {
	(*tw.contents)[tw.cy][tw.cx] = field
}

// Get will retrieve the contents of the current field
func (tw *TableWidget) Get() string {
	return (*tw.contents)[tw.cy][tw.cx]
}

// Add will add a string to the current field
func (tw *TableWidget) Add(s string) {
	(*tw.contents)[tw.cy][tw.cx] += s
}

// TrimAll will trim the leading and trailing spaces from all fields in this table
func (tw *TableWidget) TrimAll() {
	for y := 0; y < len(*tw.contents); y++ {
		for x := 0; x < len((*tw.contents)[y]); x++ {
			(*tw.contents)[y][x] = strings.TrimSpace((*tw.contents)[y][x])
		}
	}
}
