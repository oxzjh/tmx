package tmx

import (
	"encoding/xml"
	"os"
)

type TSX struct {
	Tiles []*Tile `xml:"tile"`
}

type Tile struct {
	Id         int         `xml:"id,attr"`
	Type       string      `xml:"type,attr"`
	Properties []*Property `xml:"properties>property"`
}

type Property struct {
	Name  string `xml:"name,attr"`
	Type  string `xml:"type,attr"`
	Value string `xml:"value,attr"`
}

func OpenTSX(file string) (*TSX, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var tsx TSX
	return &tsx, xml.NewDecoder(f).Decode(&tsx)
}
