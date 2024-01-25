package main

import (
	"bytes"
	"fmt"

	"html/template"

	"example.com/formbuilder/types"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/jung-kurt/gofpdf"
)

func generateValuesToPDF(form types.Form) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set font (replace "Arial" with a font available on your system)
	pdf.SetFont("Helvetica", "", 12)

	// Generate PDF using form data
	// For simplicity, just print the field names and captions in this example.
	for _, field := range form.Fields {
		pdf.Cell(0, 10, fmt.Sprintf("Field: %s, Caption: %s", field.Name, field.Caption))
		pdf.Ln(10)
	}

	// Output PDF to a file (replace with your desired output method)
	err := pdf.OutputFileAndClose("confirmation.pdf")
	if err != nil {
		fmt.Println("Error generating PDF:", err)
	}
}

func printContent(form types.Form) {

	// Print extracted values
	fmt.Printf("Program Language Field:\nName: %s\nType: %s\nOptional: %s\nFieldType: %s\nCaption: %s\n",
		form.Fields[0].Name, form.Fields[0].Type, form.Fields[0].Optional, form.Fields[0].FieldType, form.Fields[0].Caption)

	fmt.Println("\nLabels:")
	for _, label := range form.Fields[0].Labels {
		fmt.Printf("Name: %s, Value: %s\n", label.Name, label.Value)
	}

	fmt.Printf("\nExperience Section:\nName: %s\nOptional: %s\nTitle: %s\n",
		form.Sections[0].Name, form.Sections[0].Optional, form.Sections[0].Title)

	fmt.Println("\nContents:")
	for _, field := range form.Sections[0].Contents.Fields {
		fmt.Printf("Name: %s\nType: %s\nOptional: %s\nFieldType: %s\nCaption: %s\n",
			field.Name, field.Type, field.Optional, field.FieldType, field.Caption)
	}
}

func generatePDF(form types.Form) {
	// Generate HTML content using a template
	htmlContent := generateHTML(form)

	// Convert HTML to PDF using go-wkhtmltopdf
	pdfg, _ := wkhtmltopdf.NewPDFGenerator()
	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader([]byte(htmlContent))))

	err := pdfg.Create()
	if err != nil {
		fmt.Println("Error generating PDF:", err)
		return
	}

	// Write PDF to a file (replace with your desired output method)
	err = pdfg.WriteFile("confirmation.pdf")
	if err != nil {
		fmt.Println("Error writing PDF file:", err)
	}
}

func generateHTML(form types.Form) string {
	// Create an HTML template
	tmpl := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Form Confirmation</title>
		</head>
		<body>
			<h1>Form Confirmation</h1>
			{{range .Fields}}
				<p><strong>{{.Caption}}:</strong></p>
				<ul>
					{{range .Labels}}
						<li><strong>{{.Name}}:</strong> {{.Value}}</li>
					{{end}}
				</ul>
			{{end}}

			{{range .Sections}}
				<h2>{{.Title}}</h2>
				{{range .Contents.Fields}}
					<p><strong>{{.Caption}}:</strong></p>
					<ul>
						{{range .Labels}}
							<li><strong>{{.Name}}:</strong> {{.Value}}</li>
						{{end}}
					</ul>
				{{end}}
			{{end}}
		</body>
		</html>
	`

	// Parse the template
	t, err := template.New("confirmation").Parse(tmpl)
	if err != nil {
		fmt.Println("Error parsing HTML template:", err)
		return ""
	}

	// Execute the template with form data
	var htmlContent bytes.Buffer
	err = t.Execute(&htmlContent, map[string]interface{}{
		"Fields":   form.Fields,
		"Sections": form.Sections,
	})
	if err != nil {
		fmt.Println("Error executing HTML template:", err)
		return ""
	}

	return htmlContent.String()
}

func generateHTML1(form types.Form) string {
	// Create an HTML template
	tmpl := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Form Confirmation</title>
		</head>
		<body>
			<h1>Form Confirmation</h1>
			{{range .Fields}}
				<p><strong>{{.Caption}}:</strong></p>
				<ul>
					{{range .Labels}}
						<li><strong>{{.Name}}:</strong> {{.Value}}</li>
					{{end}}
				</ul>
			{{end}}
		</body>
		</html>
	`

	// Parse the template
	t, err := template.New("confirmation").Parse(tmpl)
	if err != nil {
		fmt.Println("Error parsing HTML template:", err)
		return ""
	}

	// Execute the template with form data
	var htmlContent bytes.Buffer
	err = t.Execute(&htmlContent, map[string]interface{}{
		"Fields": form.Fields,
	})
	if err != nil {
		fmt.Println("Error executing HTML template:", err)
		return ""
	}

	return htmlContent.String()
}

func main() {
	xmlData := `
	<Form>
		<Field Name="program_language" Type="Enumeration(A,B,C)" Optional="False" FieldType="Select">
			<Caption>Pick your programming language</Caption>
			<Labels>
				<Label Name="A">A(+)</Label>
				<Label Name="B">B</Label>
				<Label Name="C">C (All flavors except C#)</Label>
			</Labels>
		</Field>
		<Section Name="experience" Optional="False">
			<Title>Regarding your experience</Title>
			<Contents>
				<Field Name="other" Type="Text([0,200],Lines:4)" Optional="True" FieldType="TextBox">
					<Caption>Other programming experiences</Caption>
				</Field>
				<Field Name="code_repos" Type="File" Optional="True" FieldType="File">
					<Caption>Upload your code repo's in ZIP.</Caption>
				</Field>
			</Contents>
		</Section>
	</Form>`

	form, err := types.ParseXML(xmlData)
	if err != nil {
		fmt.Println("Error parsing XML:", err)
		return
	}
	//printContent(form)

	generatePDF(form)
}
