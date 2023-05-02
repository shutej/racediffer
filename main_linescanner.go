package main

import (
	"bufio"
	"encoding/json"
	"io"
	"sync"
	"time"

	"github.com/alecthomas/participle/v2"
)

type Datum struct {
	Value float64 `@Float`
}

func (d Datum) Int() int {
	return int(d.Value)
}

func (d Datum) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(d.Value))
}

type Line struct {
	Time time.Time
	ID   Datum   `"[" @@ "]" ":"`
	Data []Datum `( @@ "," )+`
}

var parser = participle.MustBuild[Line]()

type lineScanner struct {
	mutex sync.RWMutex
	clock clock
}

func newLineScanner() *lineScanner {
	return &lineScanner{
		clock: defaultClockInstance,
	}
}

func (ls *lineScanner) scan(s io.ReadWriter, ch chan Line) error {
	if _, err := s.Write([]byte("viewLog\r\n")); err != nil {
		return err
	}

	scanner := bufio.NewScanner(s)
	for scanner.Scan() {
		if l := ls.parse(scanner.Text()); l != nil {
			ch <- *l
		}
	}
	return scanner.Err()
}

func (ls *lineScanner) parse(s string) *Line {
	l, err := parser.ParseString("", s)
	l.Time = ls.clock.Time()
	if err == nil {
		return l
	}
	return nil
}
