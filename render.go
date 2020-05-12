package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// SavePDF can save the text as a PDF. It's pretty experimental.
func (e *Editor) SavePDF(title, filename string) error {

	// Check if the file exists
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists", filename)
	}

	// Build a large string with the document contents
	var sb strings.Builder
	for i := 0; i < e.Len(); i++ {
		// Expand tabs for each line
		sb.WriteString(strings.Replace(e.Line(LineIndex(i)), "\t", strings.Repeat(" ", e.spacesPerTab), -1) + "\n")
	}
	contents := sb.String()

	// Create a timestamp for the current date, using the "2006-01-02" format
	timestamp := time.Now().Format("2006-01-02")

	// Use A4 and Unicode
	pdf := gofpdf.New("P", "mm", "A4", "")
	tr := pdf.UnicodeTranslatorFromDescriptor("") // "" defaults to "cp1252"

	pdf.SetTopMargin(30)

	// Top text
	pdf.SetHeaderFunc(func() {
		pdf.SetY(5)
		pdf.SetFont("Helvetica", "", 6)
		// Top right corner
		pdf.CellFormat(0, 0, timestamp, "", 0, "R", false, 0, "")
	})

	// Bottom text
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Helvetica", "", 6)
		// Bottom center
		pdf.CellFormat(0, 10, fmt.Sprintf("%d", pdf.PageNo()), "", 0, "C", false, 0, "")
	})

	pdf.AddPage()
	pdf.SetY(20)
	ht := pdf.PointConvert(8.0)

	// Header
	pdf.SetFont("Courier", "B", 12)
	pdf.MultiCell(190, ht, tr(title+"\n\n"), "", "L", false)
	pdf.Ln(ht)

	// Body
	pdf.SetFont("Courier", "", 6)
	pdf.MultiCell(190, ht, tr(contents+"\n"), "", "L", false)
	pdf.Ln(ht)

	// Save to file
	return pdf.OutputFileAndClose(filename)
}
