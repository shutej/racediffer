package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	simulate = flag.String("simulate", "", "path name to read a previously captured log from RaceCapture")
	output   = flag.String("output", "", "path name to write *.jsonl file")
)

var ls = newLineScanner()

type Event struct {
	Time  time.Time
	Event int
}

type model struct {
	clock clock
	count int
	ch    chan Event
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeySpace {
			m.ch <- Event{Time: m.clock.Time(), Event: m.count}
			m.count++
			return m, nil
		}
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *model) View() string {
	return fmt.Sprintf("Press [SPACE] to send event #%d or [CTRL+C] to quit.", m.count)
}

func main() {
	flag.Parse()
	if *output == "" {
		log.Fatal("The --output flag is required.")
	}

	var s io.ReadWriteCloser
	var err error
	if *simulate == "" {
		s, err = openSerial()
	} else {
		s, err = openSimulator(*simulate)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	chEvent := make(chan Event)
	chLine := make(chan Line)

	go func() {
		defer close(chLine)
		ls.scan(s, chLine)
	}()

	done := make(chan struct{})
	f, err := os.Create(*output)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	go func() {
		defer close(done)
	done:
		for {
			select {
			case line := <-chLine:
				b, err := json.Marshal(line)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Fprintf(f, "%s\n", b)
			case event, ok := <-chEvent:
				if !ok {
					break done
				}
				b, err := json.Marshal(event)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Fprintf(f, "%s\n", b)
			}
		}
	}()

	if _, err := tea.NewProgram(&model{clock: defaultClockInstance, ch: chEvent}).Run(); err != nil {
		log.Fatal(err)
	}
	close(chEvent)
	<-done
}
