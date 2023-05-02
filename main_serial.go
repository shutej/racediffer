package main

import (
	"flag"
	"io"

	"github.com/tarm/serial"
)

var (
	serialDevice = flag.String("device", "/dev/cu.usbmodem00000000011C1", "the path to the serial device")
	serialBaud   = flag.Int("baud", 115200, "the baud of the serial device")
)

func openSerial() (io.ReadWriteCloser, error) {
	c := &serial.Config{Name: *serialDevice, Baud: *serialBaud}
	return serial.OpenPort(c)
}
