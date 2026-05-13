package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/serverhorror/tiffcheck/pkg/tiffcheck"
)

func TestRun_HelpFlag(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	code := run([]string{"--help"}, &out, &errOut)
	if code != 0 {
		t.Fatalf("expected code 0, got %d", code)
	}
	if errOut.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", errOut.String())
	}
	if !strings.Contains(out.String(), "Usage: tiffcheck") {
		t.Fatalf("expected usage output, got %q", out.String())
	}
}

func TestRun_VersionFlag(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	code := run([]string{"--version"}, &out, &errOut)
	if code != 0 {
		t.Fatalf("expected code 0, got %d", code)
	}
	if errOut.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", errOut.String())
	}
	stdout := out.String()
	if !strings.Contains(stdout, "tiffcheck") {
		t.Fatalf("expected binary name in output, got %q", stdout)
	}
	if !strings.Contains(stdout, "revision:") {
		t.Fatalf("expected revision in output, got %q", stdout)
	}
	if !strings.Contains(stdout, "time:") {
		t.Fatalf("expected time in output, got %q", stdout)
	}
}

func TestRun_InvalidFlag(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	code := run([]string{"--bogus"}, &out, &errOut)
	if code != 2 {
		t.Fatalf("expected code 2, got %d", code)
	}
	if out.Len() != 0 {
		t.Fatalf("expected empty stdout, got %q", out.String())
	}
	stderr := errOut.String()
	if !strings.Contains(stderr, "Error:") {
		t.Fatalf("expected parse error, got %q", stderr)
	}
	if !strings.Contains(stderr, "Usage: tiffcheck") {
		t.Fatalf("expected usage text, got %q", stderr)
	}
}

func TestRun_MissingFileArgument(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	code := run([]string{}, &out, &errOut)
	if code != 2 {
		t.Fatalf("expected code 2, got %d", code)
	}
	if out.Len() != 0 {
		t.Fatalf("expected empty stdout, got %q", out.String())
	}
	if !strings.Contains(errOut.String(), "Usage: tiffcheck") {
		t.Fatalf("expected usage text, got %q", errOut.String())
	}
}

func TestRun_OpenFailure(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	code := run([]string{filepath.Join(t.TempDir(), "missing.tif")}, &out, &errOut)
	if code != 1 {
		t.Fatalf("expected code 1, got %d", code)
	}
	if out.Len() != 0 {
		t.Fatalf("expected empty stdout, got %q", out.String())
	}
	if !strings.Contains(errOut.String(), "Error:") {
		t.Fatalf("expected open error, got %q", errOut.String())
	}
}

func TestRun_CheckFailure(t *testing.T) {
	tempFile := filepath.Join(t.TempDir(), "any.tif")
	if err := os.WriteFile(tempFile, []byte("content"), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	original := check
	check = func(tiffcheck.Reader, ...tiffcheck.Option) (tiffcheck.Report, error) {
		return tiffcheck.Report{}, errors.New("boom")
	}
	defer func() { check = original }()

	var out bytes.Buffer
	var errOut bytes.Buffer

	code := run([]string{tempFile}, &out, &errOut)
	if code != 1 {
		t.Fatalf("expected code 1, got %d", code)
	}
	if out.Len() != 0 {
		t.Fatalf("expected empty stdout, got %q", out.String())
	}
	if !strings.Contains(errOut.String(), "Error: boom") {
		t.Fatalf("expected check error, got %q", errOut.String())
	}
}

func TestRun_Success(t *testing.T) {
	tempFile := filepath.Join(t.TempDir(), "any.tif")
	if err := os.WriteFile(tempFile, []byte("content"), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	original := check
	check = func(tiffcheck.Reader, ...tiffcheck.Option) (tiffcheck.Report, error) {
		return tiffcheck.Report{
			Tiled:     true,
			TileWidth: 256,
			GeoTIFF:   true,
		}, nil
	}
	defer func() { check = original }()

	var out bytes.Buffer
	var errOut bytes.Buffer

	code := run([]string{tempFile}, &out, &errOut)
	if code != 0 {
		t.Fatalf("expected code 0, got %d", code)
	}
	if errOut.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", errOut.String())
	}

	stdout := out.String()
	if !strings.Contains(stdout, "Status:  TILED (Indexed for partial loading)") {
		t.Fatalf("expected tiled output, got %q", stdout)
	}
	if !strings.Contains(stdout, "Tile Dim: 256 pixels wide") {
		t.Fatalf("expected tile dim output, got %q", stdout)
	}
	if !strings.Contains(stdout, "Metadata: GeoTIFF tags detected") {
		t.Fatalf("expected geotiff output, got %q", stdout)
	}
}
