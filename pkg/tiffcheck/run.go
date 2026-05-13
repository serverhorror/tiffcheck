package tiffcheck

import (
	"errors"
	"io"
)

type Option func(*config)

type Reader interface {
	io.ReaderAt
	io.ReadSeeker
}

type config struct {
	analyzer Analyzer
}

func defaultConfig() config {
	return config{
		analyzer: NewGoogleTIFFAnalyzer(),
	}
}

func WithAnalyzer(analyzer Analyzer) Option {
	return func(cfg *config) {
		if analyzer != nil {
			cfg.analyzer = analyzer
		}
	}
}

func Check(r Reader, opts ...Option) (Report, error) {
	if r == nil {
		return Report{}, errors.New("reader is required")
	}

	cfg := defaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	report, err := cfg.analyzer.Analyze(r)
	if err != nil {
		return Report{}, err
	}
	return report, nil
}
