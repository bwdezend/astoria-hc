package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/vemo-france/max31865"
)

var acc accessory.Thermostat

func init() {
	info := accessory.Info{
		Name: "Astoria Boiler Thermostat",
	}
	temp := 10.0
	min := 10.0
	max := 130.0
	steps := 0.1
	acc = *accessory.NewThermostat(info, temp, min, max, steps)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var path = flag.String("path", "/dev/shm", "path to current_temp")
var gpio = flag.Bool("gpio", true, "load gpio code")
var p = 3.0

func getCurrentTemp() {
	sensor := max31865.Create("8", "9", "10", "11")
	var boilerTemperature float64
	for {
		boilerTemperature = float64(sensor.ReadTemperature(100, 430))
		acc.Thermostat.CurrentTemperature.SetValue(boilerTemperature)
		currentTemperature.Set(boilerTemperature)
		time.Sleep(500 * time.Millisecond)
	}

}

//func getCurrentTemp() {
//	// Reads the boiler current temp from disk, as output by the python program
//	// that is translating the RTD amplifier signal to degrees C.
//
//	// TODO: decode the RTD directly inside this function
//
//	var newTemp float64 = 0
//
//	for {
//		f, err := os.Open(filepath.Join(*path, "current_temp"))
//		check(err)
//
//		sc := bufio.NewScanner(f)
//		for sc.Scan() {
//			newTemp, err = strconv.ParseFloat(sc.Text(), 64)
//		}
//		f.Close()
//		check(err)
//
//		acc.Thermostat.CurrentTemperature.SetValue(newTemp)
//		currentTemperature.Set(newTemp)
//		time.Sleep(500 * time.Millisecond)
//	}
//
//}

func setTargetTemp(setTemp float64) {
	// Adjust the setpoint for the PID loop and updates the homekit interfaces

	if setTemp > 124.0 {
		setTemp = 124.0
	}

	log.Printf("setting setpoint to %.2f", setTemp)
	acc.Thermostat.TargetTemperature.SetValue(setTemp)
	setpointTemperature.Set(setTemp)
}

func heaterWindow(windowSize float64, enabledTime float64) {
	// Take two inputs - the duration of the cycle and the proportion of the cycle
	// and turn the heating element on for that percentage of the cycle. If the
	// windowSize is 15.0 and the enabledTime is 0.7, this turns the heating element
	// on for 10.5 seconds and off for 4.5 seconds before returning

	var disabledTime float64 = 0
	enabledTime = windowSize * enabledTime * 1000
	disabledTime = windowSize*1000 - enabledTime

	if enabledTime > 1000.0 {
		enabledTime = 1000.0
	}
	if enabledTime > 0 {
		if *gpio {
			heaterControl(true)
		}
		time.Sleep(time.Duration(enabledTime) * time.Millisecond)
		secondsActive.Add(enabledTime / 1000)
	}
	if *gpio {
		heaterControl(false)
	}
	time.Sleep(time.Duration(disabledTime) * time.Millisecond)
}

func temperatureProportional() {
	// Dead simple proportional control loop. Take the difference in setpoint
	// and current temp, multiply by the gain, and use the result to control
	// the duty cycle on the boiler, represented as a float between 0.0 and 1.0

	for {
		current := acc.Thermostat.CurrentTemperature.GetValue()
		setpoint := acc.Thermostat.TargetTemperature.GetValue()
		error := (setpoint - current) * p * 0.1

		if error > 1.0 {
			error = 1.0
		}

		if error < 0.01 {
			error = 0.0
		}

		// log.Printf("Duty Cycle: %.2f, Current: %.2f, Setpoint: %.2f\n", error, current, setpoint)

		heaterWindow(1.0, error)

	}
}

func temperatureErrorDetection() {
	// This is a error handling function to turn the SSR off if the boiler doesn't
	// appear to be warming as a result of input. The most common reason for this
	// is that the power switch on the machine is turned off, but the timer still
	// fired. Less common reasons would be a heater element malfunction, a boiler
	// rupture, etc.

	lastTemp := 0.0
	countdown := 2
	for {
		setpoint := acc.Thermostat.TargetTemperature.GetValue()
		time.Sleep(10 * time.Second)

		current := acc.Thermostat.CurrentTemperature.GetValue()
		error := (setpoint - current)

		if error > 5.0 {
			if lastTemp >= current {
				if countdown > 0 {
					log.Printf("Warning: boiler does not appear to be heating (%.2f gt %.2f). Countdown: %d", lastTemp, current, countdown)
					countdown--
				} else if countdown == 0 {
					log.Printf("Error: boiler does not appear to be heating (%.2f gt %.2f). Changing setpoint to 0.0", lastTemp, current)
					setTargetTemp(0.0)
					countdown = 2
				}
			}
		} else if countdown < 2 {
			log.Printf("Warning: boiler does not appear to be heating (%.2f gt %.2f). Countdown: %d", lastTemp, current, countdown)
			countdown = 2
		}
		lastTemp = current
	}
}

func main() {
	flag.Parse()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()

	config := hc.Config{Pin: "00102003"}
	t, err := hc.NewIPTransport(config, acc.Accessory)
	check(err)
	hc.OnTermination(func() {
		t.Stop()
	})

	go prometheusMetrics()

	// Picks up the current temperature every second from *path/current_temp
	go getCurrentTemp()

	// Run the PID loop itself
	go temperatureProportional()

	// Enable the "boiler isn't heating" error handling
	go temperatureErrorDetection()

	// Runs the REST-y interface for adjusting the setpoint and causing the program to exit.
	go setpointHandler()

	// If the -gpio flag is set to true load the Raspberry Pi gpio handling code.
	if *gpio {
		log.Println("rpi-gpio enabled - enabling button on pin 26")
		go hwButton()
	}

	// Enable the homekit code to change the setpoint
	acc.Thermostat.TargetTemperature.OnValueRemoteUpdate(func(temp float64) {
		setTargetTemp(temp)
	})

	t.Start()

}
