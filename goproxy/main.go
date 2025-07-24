package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/goproxy/goproxy"
)

func main() {
	// Set up the handler with desired options
	opts := slog.HandlerOptions{
		AddSource: true,            // Include source file and line number
		Level:     slog.LevelDebug, // Set the minimum level to DEBUG
	}
	handler := slog.NewTextHandler(os.Stdout, &opts)
	logger := slog.New(handler)
	http.ListenAndServe("0.0.0.0:80", &goproxy.Goproxy{
		Fetcher: &fetcher{
			Fetcher: &goproxy.GoFetcher{},
			logger:  logger.WithGroup("Fetcher"),
		},
		Logger: logger,
	})
}

type fetcher struct {
	goproxy.Fetcher

	logger *slog.Logger
}

func (f *fetcher) Query(ctx context.Context, path, query string) (version string, time time.Time, err error) {
	f.logger.Debug("Querying", "path", path, "query", query)
	return f.Fetcher.Query(ctx, path, query)
}

func (f *fetcher) List(ctx context.Context, path string) (versions []string, err error) {
	f.logger.Debug("Listing versions", "path", path)
	return f.Fetcher.List(ctx, path)
}

func (f *fetcher) Download(ctx context.Context, path, version string) (info, mod, zip io.ReadSeekCloser, err error) {
	f.logger.Debug("Downloading", "path", path, "version", version)
	if path == "go.opentelemetry.io/auto" && version == "v0.22.1" {
		f.logger.Debug("Serving local files for go.opentelemetry.io/auto v0.22.1")
		i, err := os.Open("./info.txt")
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to open info.txt: %w", err)
		}
		m, err := os.Open("./mod.txt")
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to open mod.txt: %w", err)
		}
		z, err := os.Open("./auto-v0.22.1.zip")
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to open auto-v0.22.1.zip: %w", err)
		}
		return i, m, z, nil

		/*
			i, m, z, err := f.Fetcher.Download(ctx, path, version)
			f.logger.Debug("Special case for go.opentelemetry.io/auto", "path", path, "version", version, "info", i, "mod", m, "zip", z, "error", err)
			infoF, e := os.OpenFile("info.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			err = errors.Join(err, e)
			if err == nil {
				i = &teeReadSeekCloser{ReadSeekCloser: i, f: infoF}
			}
			modF, e := os.OpenFile("mod.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			err = errors.Join(err, e)
			if err == nil {
				m = &teeReadSeekCloser{ReadSeekCloser: m, f: modF}
			}
			zipF, e := os.OpenFile("auto-v0.22.1.zip", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			err = errors.Join(err, e)
			if err == nil {
				z = &teeReadSeekCloser{ReadSeekCloser: z, f: zipF}
			}
			return i, m, z, err
		*/
	}
	return f.Fetcher.Download(ctx, path, version)
}

type teeReadSeekCloser struct {
	io.ReadSeekCloser

	f *os.File
}

func (t *teeReadSeekCloser) Read(p []byte) (n int, err error) {
	tee := io.TeeReader(t.ReadSeekCloser, t.f)
	return tee.Read(p)
}

func (t *teeReadSeekCloser) Close() error {
	err := t.f.Close()
	return errors.Join(err, t.ReadSeekCloser.Close())
}
