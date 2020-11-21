package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/bwdezend/astoria-hc/internal/core"
	"github.com/bwdezend/astoria-hc/internal/telemetry"
)

var acc accessory.Thermostat
var deviceName = flag.String("name", "Astoria Boiler Thermostat", "Device name")
var dbPath = flag.String("db", "./db", "Database path")
var enableMetics = flag.Bool("metrics", true, "Enable prometheus metrics")
var prometheusPort = flag.Int("promPort", 2112, "Port to reigster /metrics handler on")
var temperatureMinimim = flag.Float64("minTemp", 10.0, "Minimum temperature value")
var temperatureMaximum = flag.Float64("maxTemp", 130.0, "Maximum temperature value")
var temperatureStepSize = flag.Float64("stepTemp", 0.1, "Temperature setting step size")
var homekitPin = flag.String("pin", "00102003", "Homekit Pairing PIN")

var gpio = flag.Bool("gpio", true, "load gpio code")

func init() {
	flag.Parse()
	info := accessory.Info{
		Name: *deviceName,
	}

	acc = *accessory.NewThermostat(info, 10.0, *temperatureMinimim, *temperatureMaximum, *temperatureStepSize)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var p = 3.0

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()

	config := hc.Config{Pin: *homekitPin}
	t, err := hc.NewIPTransport(config, acc.Accessory)
	check(err)
	hc.OnTermination(func() {
		t.Stop()
	})

	if *enableMetics {
		go telemetry.PrometheusMetrics(*prometheusPort)
	}

	// Picks up the current temperature every second from *path/current_temp
	go core.GetCurrentTemp(acc)

	// Run the PID loop itself
	go core.TemperatureProportional(acc, *gpio)

	// Enable the "boiler isn't heating" error handling
	go core.TemperatureErrorDetection(acc)

	// Runs the REST-y interface for adjusting the setpoint and causing the program to exit.
	//go telemetry.SetpointHandler()

	// If the -gpio flag is set to true load the Raspberry Pi gpio handling code.
	if *gpio {
		log.Println("rpi-gpio enabled - enabling power button on pin 26")
		go core.PowerButton(acc)
	}

	// Enable the homekit code to change the setpoint
	acc.Thermostat.TargetTemperature.OnValueRemoteUpdate(func(temp float64) {
		core.SetTargetTemp(acc, temp)
	})

	t.Start()

}
