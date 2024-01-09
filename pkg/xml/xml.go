package xml

import "encoding/xml"

type Marshaller interface {
	MarshalIndent(v interface{}) ([]byte, error)
}

type DefaultXMLMarshaller struct{}

func (m DefaultXMLMarshaller) MarshalIndent(v interface{}) ([]byte, error) {
	return xml.MarshalIndent(v, "", " ")
}
