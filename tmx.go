package tmx

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"io"
	"os"
	"strings"
)

type TMX struct {
	Width    int        `xml:"width,attr"`
	Height   int        `xml:"height,attr"`
	Tilesets []*Tileset `xml:"tileset"`
	Layers   []*Layer   `xml:"layer"`
}

type Tileset struct {
	Firstgid int    `xml:"firstgid,attr"`
	Source   string `xml:"source,attr"`
}

type Layer struct {
	Name string `xml:"name,attr"`
	Data *Data  `xml:"data"`
}

type Data struct {
	Encoding    string `xml:"encoding,attr"`
	Compression string `xml:"compression,attr"`
	Raw         string `xml:",innerxml"`
	Gids        []int
}

func (d *Data) decode(n int) error {
	d.Raw = strings.TrimSpace(d.Raw)
	switch d.Encoding {
	case "csv":
		d.Gids = make([]int, 0, n)
		gid := 0
		r := strings.NewReader(d.Raw)
		for {
			b, err := r.ReadByte()
			if err != nil {
				return nil
			}
			if b >= '0' && b <= '9' {
				gid *= 10
				gid += int(b - '0')
			} else if b == ',' {
				d.Gids = append(d.Gids, gid)
				gid = 0
			}
		}
	case "base64":
		b, err := base64.StdEncoding.DecodeString(d.Raw)
		if err != nil {
			return err
		}
		var r io.Reader
		switch d.Compression {
		case "":
		case "gzip":
			if r, err = gzip.NewReader(bytes.NewReader(b)); err != nil {
				return err
			}
		case "zlib":
			if r, err = zlib.NewReader(bytes.NewReader(b)); err != nil {
				return err
			}
		default:
			return errors.New("unsupported compression: " + d.Compression)
		}
		if r != nil {
			if b, err = io.ReadAll(r); err != nil {
				return err
			}
		}
		d.Gids = make([]int, n)
		for i := 0; i < n; i++ {
			o := i << 2
			d.Gids[i] = int(b[o]) | int(b[o+1])<<8 | int(b[o+2])<<16 | int(b[o+3])<<24
		}
		return nil
	default:
		return errors.New("unsupported encoding: " + d.Encoding)
	}
}

func Open(file string) (*TMX, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var tmx TMX
	if err = xml.NewDecoder(f).Decode(&tmx); err != nil {
		return nil, err
	}
	n := tmx.Width * tmx.Height
	for _, layer := range tmx.Layers {
		if err = layer.Data.decode(n); err != nil {
			return nil, err
		}
	}
	return &tmx, nil
}
