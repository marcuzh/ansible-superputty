package superputty

import (
	"encoding/xml"
	"fmt"

	xml2 "github.com/marcuzh/ansible-superputty/pkg/xml"
)

type SessionData struct {
	XMLName xml.Name `xml:"ArrayOfSessionData"`
	XSD     string   `xml:"xmlns:xsd,attr"`
	XSI     string   `xml:"xmlns:xsi,attr"`
}

func (s *SessionData) ToXML(marshaller xml2.Marshaller) (string, error) {
	xmlData, err := marshaller.MarshalIndent(s)

	if err != nil {
		return "", fmt.Errorf("session data marshal error: %w", err)
	}
	return xml.Header + string(xmlData), nil
}
