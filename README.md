# racediffer

This is a tool intended to work with [RaceCapture's CAN Bus Data
Logger](https://www.autosportlabs.com/use-racecapture-as-a-general-purpose-can-bus-logger-and-reverse-engineer-your-cars-data/).
It eliminates the need for tools like Hyperterminal and Minicom and directly
decodes the output of the default RaceCapture script, logging to a *.jsonl
file.

It also allows you to hit the space bar to emit a numbered event to the log
file.  This is particularly handy when you're trying to reverse engineer the
CAN bus by following a script and demarcating specific events that happen in
the car, for example:

    * Get in the car and turn on accessories.  Wait 30 seconds.
    * EVENT 1
    * Throttle up
    * EVENT 2
    * EVENT 3
    * Throttle down
    * EVENT 4

When reverse engineering, you could now compare IDs looking for ones that are
flat before EVENT 1, rising between EVENT 1 and EVENT 2, decreasing between
EVENT 3 and EVENT 4.

The only required flag specifies where to write the output, as the terminal UI
is written to STDOUT and log messages are written to STDERR:

    * --output: path name to write *.jsonl file

You can capture output directly from the serial port:

    * --device: the path to the serial device (default /dev/cu.usbmodem00000000011C1)
    * --baud: the baud of the serial device (default 115200)

If you captured a log from a RaceCapture previously, you can replay it instead
of connecting through the serial port, and control the rate that the log is
read:

    * --simulate: path name to read a previously captured log from RaceCapture
    * --simulate_rate: the rate at which the simulated rate limited reader bucket refills (default 8192)
    * --simulate_capacity: the capacity of the simulated rate limited reader bucket (default 8192)
