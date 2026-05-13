package tiffcheck

import (
	"errors"

	"github.com/google/tiff"
)

var ErrNoIFDs = errors.New("no IFDs found")

type Report struct {
	Tiled     bool
	TileWidth any
	GeoTIFF   bool
}

type Analyzer interface {
	Analyze(Reader) (Report, error)
}

type GoogleTIFFAnalyzer struct{}

func NewGoogleTIFFAnalyzer() Analyzer {
	return GoogleTIFFAnalyzer{}
}

func (GoogleTIFFAnalyzer) Analyze(r Reader) (Report, error) {
	parsed, err := tiff.Parse(r, tiff.DefaultTagSpace, tiff.DefaultFieldTypeSpace)
	if err != nil {
		return Report{}, err
	}

	ifds := parsed.IFDs()
	if len(ifds) == 0 {
		return Report{}, ErrNoIFDs
	}

	ifd := ifds[0]
	tileWidthTag := ifd.GetField(tagTileWidth)
	report := Report{
		Tiled: tileWidthTag != nil,
		GeoTIFF: ifd.GetField(tagModelPixelScale) != nil ||
			ifd.GetField(tagGeoKeyDirectory) != nil,
	}
	if tileWidthTag != nil {
		report.TileWidth = tileWidthTag.Value()
	}

	return report, nil
}
