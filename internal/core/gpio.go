package core

import (
	"log"
	"time"

	"github.com/brutella/hc/accessory"
	"github.com/bwdezend/astoria-hc/internal/telemetry"
	"github.com/stianeikeland/go-rpio"
)

func init() {
	rpio.Open()
}

// PowerButton sets up the button hooked to pin 26
// on the Rasberry Pi to act as a high/low control
// for the system.
func PowerButton(acc accessory.Thermostat) {
	high := 124.0
	low := 30.0

	button := rpio.Pin(26)
	button.Input()
	button.PullUp()

	for {
		if button.Read() == 0 {
			log.Println("button pressed")
			if acc.Thermostat.TargetTemperature.GetValue() >= (high - 10.0) {
				SetTargetTemp(acc, low)
			} else {
				SetTargetTemp(acc, high)
			}
			time.Sleep(500 * time.Millisecond)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// HeaterControl turns on and off voltate to the
// pin hooked to the Solid State Relay.
func HeaterControl(on bool) {
	pin := rpio.Pin(14)
	pin.Output()

	telemetry.RelayActivations.Inc()
	if on {
		pin.High()
	} else {
		pin.Low()
	}
}
