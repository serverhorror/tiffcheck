package tiffcheck_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/serverhorror/tiffcheck/pkg/tiffcheck"
)

func TestCheck_UsesAnalyzerOption(t *testing.T) {
	want := tiffcheck.Report{
		Tiled:     true,
		TileWidth: 256,
		GeoTIFF:   true,
	}

	got, err := tiffcheck.Check(bytes.NewReader([]byte("content")), tiffcheck.WithAnalyzer(fakeAnalyzer{
		report: tiffcheck.Report{
			Tiled:     true,
			TileWidth: 256,
			GeoTIFF:   true,
		},
	}))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got != want {
		t.Fatalf("expected report %+v, got %+v", want, got)
	}
}

func TestCheck_AnalyzerError(t *testing.T) {
	_, err := tiffcheck.Check(bytes.NewReader([]byte("content")), tiffcheck.WithAnalyzer(fakeAnalyzer{err: errors.New("parse failed")}))
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

func TestCheck_NilReaderError(t *testing.T) {
	_, err := tiffcheck.Check(nil, tiffcheck.WithAnalyzer(fakeAnalyzer{}))
	if err == nil {
		t.Fatal("expected reader error, got nil")
	}
}

func TestWithAnalyzer_NilKeepsDefault(t *testing.T) {
	opt := tiffcheck.WithAnalyzer(nil)
	if opt == nil {
		t.Fatal("expected option function, got nil")
	}
}

type fakeAnalyzer struct {
	report tiffcheck.Report
	err    error
}

func (f fakeAnalyzer) Analyze(tiffcheck.Reader) (tiffcheck.Report, error) {
	if f.err != nil {
		return tiffcheck.Report{}, f.err
	}
	return f.report, nil
}
