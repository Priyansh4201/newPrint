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

	// 1. Initialize generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}

	// 2. Open file
	f, err := os.Open(inputHTML)
	if err != nil {
		log.Fatalf("Error opening HTML: %v", err)
	}
	defer f.Close()

	// 3. Create Page
	page := wkhtmltopdf.NewPageReader(f)

	// --- CRITICAL SETTINGS TO FIX YOUR ERRORS ---
	
	// Ignore missing images and network errors
	page.LoadErrorHandling.Set("ignore")
	page.LoadMediaErrorHandling.Set("ignore")
	
	// Disable local file access if you don't have the /assets/ folder on your C: drive
	// This prevents the "system cannot find the path specified" error
	page.EnableLocalFileAccess.Set(false) 

	// Give it enough time to handle the large 1.9MB file
	page.JavascriptDelay.Set(200)
	// --------------------------------------------

	pdfg.AddPage(page)

	// Global settings
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Dpi.Set(300)
	pdfg.MarginLeft.Set(0)
	pdfg.MarginRight.Set(0)
	pdfg.MarginTop.Set(0)
	pdfg.MarginBottom.Set(0)

	fmt.Printf("⏳ Converting %s (1.9MB) to PDF...\n", inputHTML)
	start := time.Now()

	// 4. Create the PDF
	err = pdfg.Create()
	
	// If it still returns an error, check if it actually managed to generate bytes.
	// Sometimes wkhtmltopdf returns an error code even if the PDF was produced.
	if err != nil {
		fmt.Printf("⚠️ Note: wkhtmltopdf reported warnings: %v\n", err)
		if len(pdfg.Bytes()) == 0 {
			log.Fatal("❌ Failed to generate any PDF content.")
		}
	}

	// 5. Save to file
	err = pdfg.WriteFile(outputPDF)
	if err != nil {
		log.Fatalf("Error saving file: %v", err)
	}

	fmt.Printf("✅ Done! File saved as %s in %v\n", outputPDF, time.Since(start))
}