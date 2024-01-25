package types

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

// Form represents the structure of the form data.
type Form struct {
	XMLName  xml.Name  `xml:"Form"`
	Fields   []Field   `xml:"Field"`
	Sections []Section `xml:"Section"`
}

// Field represents a form field.
type Field struct {
	Name      string  `xml:"Name,attr"`
	Type      string  `xml:"Type,attr"`
	Optional  string  `xml:"Optional,attr"`
	FieldType string  `xml:"FieldType,attr"`
	Caption   string  `xml:"Caption"`
	Labels    []Label `xml:"Labels>Label"`
}

// Label represents a label within a field.
type Label struct {
	Name  string `xml:"Name,attr"`
	Value string `xml:",chardata"`
}

// Section represents a form section.
type Section struct {
	Name     string   `xml:"Name,attr"`
	Optional string   `xml:"Optional,attr"`
	Title    string   `xml:"Title"`
	Contents Contents `xml:"Contents"`
}

// Contents represents the contents of a form section.
type Contents struct {
	Fields   []Field   `xml:"Field"`
	Sections []Section `xml:"Section"`
}

// ParseXML parses the XML data into a Form structure.
func ParseXML(xmlData string) (Form, error) {
	var form Form
	err := xml.Unmarshal([]byte(xmlData), &form)
	return form, err
}

// ReadXMLFromFile reads XML data from a file and parses it into a Form structure.
func ReadXMLFromFile(filePath string) (Form, error) {
	xmlData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Form{}, fmt.Errorf("error reading XML file: %v", err)
	}

	return ParseXML(string(xmlData))
}
