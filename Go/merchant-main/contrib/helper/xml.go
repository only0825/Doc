package helper

import (
	"encoding/xml"
)

func XMLMarshal(v interface{}) ([]byte, error) {

	output, err := xml.MarshalIndent(v, "  ", "    ")
	if err != nil {
		return nil, err
	}

	return []byte(xml.Header + "\n" + string(output)), nil
}

func XMLUnmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, &v)
}
