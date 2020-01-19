package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "golang.org/x/image/vector"

	"github.com/unidoc/unipdf/v3/creator"
	pdf "github.com/unidoc/unipdf/v3/model"
)

func main() {
	operation := os.Args[1]

	switch operation {
	case "meta":
		operationMeta()
	case "sign":
		operationSign(os.Args[2])
	case "verify":
		operationVerify()
	default:
		fmt.Fprintf(os.Stderr, "Unknown operation")
		os.Exit(1)
	}
}

func operationMeta() {
	fmt.Println("--RESULT--")
	fmt.Println("OK")
	// TODO: Add support for images - image/bmp image/webp image/jpeg image/tiff image/png image/gif
	fmt.Println("SUPPORTS:application/pdf")
	fmt.Println("--RESULT--")
	os.Exit(0)
}

func operationSign(inputPath string) {
	imagePath := prompt("image", "Signature file", "")

	position := prompt("position", "Signature position", imagePath)

	positionParts := strings.Split(position, ",")

	xPos, err := strconv.ParseFloat(positionParts[0], 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	yPos, err := strconv.ParseFloat(positionParts[1], 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	iwidth, err := strconv.ParseFloat(positionParts[2], 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputPath := prompt("save", "Output path", inputPath)
	// Make sure path has .pdf extension
	hasExt := strings.HasSuffix(outputPath, ".pdf")
	if !hasExt {
		outputPath = outputPath + ".pdf"
	}

	pageNumStr := positionParts[3]
	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	err = addImageToPdf(inputPath, outputPath, imagePath, pageNum, xPos, yPos, iwidth)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func operationVerify() {
	fmt.Println("--RESULT--")
	fmt.Println("UNKNOWN")
	fmt.Println("--RESULT--")
	os.Exit(0)
}

func prompt(promptType string, question string, defaultValue string) string {
	fmt.Println("--PROMPT--")
	fmt.Printf("%s\"%s\"(\"%s\")\n", promptType, question, defaultValue)
	fmt.Println("--PROMPT--")

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	if strings.TrimSpace(text) == "--PROMPT--" {
		response, _ := reader.ReadString('\n')

		// Empty response
		if strings.TrimSpace(response) == "" {
			fmt.Fprintf(os.Stderr, "No answer, aborting.")
			os.Exit(1)
			return ""
		}

		reader.ReadString('\n')
		return strings.TrimSpace(response)
	}

	fmt.Fprintf(os.Stderr, "Unexpected response to prompt '%s'.", text)
	os.Exit(1)
	return ""
}

// Inspired by the example from unipdf https://github.com/unidoc/unipdf-examples/blob/v3/image/pdf_add_image_to_page.go
func addImageToPdf(inputPath string, outputPath string, imagePath string, pageNum int, xPos float64, yPos float64, iwidth float64) error {

	c := creator.New()

	// Prepare the image.
	img, err := c.NewImageFromFile(imagePath)
	if err != nil {
		return err
	}
	img.ScaleToWidth(iwidth)
	img.SetPos(xPos, yPos)

	// Read the input pdf file.
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	// Load the pages.
	for i := 0; i < numPages; i++ {
		page, err := pdfReader.GetPage(i + 1)
		if err != nil {
			return err
		}

		// Add the page.
		err = c.AddPage(page)
		if err != nil {
			return err
		}

		// If the specified page, or -1, apply the image to the page.
		if i+1 == pageNum || pageNum == -1 {
			_ = c.Draw(img)
		}
	}

	err = c.WriteToFile(outputPath)
	return err
}
