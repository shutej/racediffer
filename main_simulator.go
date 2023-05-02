package main

import (
	"flag"
	"io"
	"os"

	"github.com/juju/ratelimit"
)

var (
	simulateRate     = flag.Float64("simulate_rate", 8192, "the rate at which the simulated rate limited reader bucket refills")
	simulateCapacity = flag.Int64("simulate_capacity", 8192, "the capacity of the simulated rate limited reader bucket")
)

type nopWriter struct {
	io.ReadCloser
	io.Writer
}

type ratelimitReadCloser struct {
	io.Reader
	io.Closer
}

func openSimulator(path string) (io.ReadWriteCloser, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	b := ratelimit.NewBucketWithRate(*simulateRate, *simulateCapacity)
	return nopWriter{
		ReadCloser: ratelimitReadCloser{
			Reader: ratelimit.Reader(r, b),
			Closer: r,
		},
		Writer: io.Discard,
	}, nil
}
