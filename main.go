package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func main() {
	inputHTML := "pdfgo.html"
	outputPDF := "output.pdf"

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(inputHTML)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	page := wkhtmltopdf.NewPageReader(f)

	// --- FIXING THE ERRORS ---
	
	// 1. Correct name for the "Intelligent Shrinking" setting
	page.DisableSmartShrinking.Set(true) 
	
	// 2. Standard Viewport for A4 at 96 DPI
	page.ViewportSize.Set("794x1122")
	
	// 3. Essential for Marathi/Hindi characters
	page.Encoding.Set("utf-8")
	
	// 4. Ignore the missing icon/signature errors
	page.LoadErrorHandling.Set("ignore")
	page.DisableJavascript.Set(true)

	pdfg.AddPage(page)

	// --- GLOBAL A4 CALIBRATION ---
	
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	
	// IMPORTANT: Set DPI to 96. 
	// At 300 DPI, your 210mm width is interpreted differently.
	// 96 DPI ensures 1 pixel in CSS = 1/96 inch (Standard Web).
	pdfg.Dpi.Set(96) 

	// Remove all engine margins
	pdfg.MarginLeft.Set(0)
	pdfg.MarginRight.Set(0)
	pdfg.MarginTop.Set(0)
	pdfg.MarginBottom.Set(0)

	fmt.Printf("⏳ Generating exact A4 PDF...")
	start := time.Now()
	err = pdfg.Create()

	if err != nil && len(pdfg.Bytes()) == 0 {
		log.Fatalf("❌ Failed: %v", err)
	}

	err = pdfg.WriteFile(outputPDF)
	fmt.Printf("✅ Done in %v\n", time.Since(start))
}