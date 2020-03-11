package main

import (
    "github.com/stianeikeland/go-rpio"
    "time"
    "log"
)

func init() {
    rpio.Open()
}

func hwButton() {
    high := 124.0
    low  := 30.0

    button := rpio.Pin(26)
    button.Input()
    button.PullUp()

    for {
        if button.Read() == 0 {
            log.Println("button pressed")
            if acc.Thermostat.TargetTemperature.GetValue() >= ( high - 10.0 ) {
                setTargetTemp(low)
            } else {
                setTargetTemp(high)
            }
            time.Sleep(500 * time.Millisecond)
         }
         time.Sleep(50 * time.Millisecond)
    }
}

func heaterControl(on bool) {
    pin := rpio.Pin(14)
    pin.Output()

    relayActivations.Inc()
    if on{
        pin.High()
    } else {
        pin.Low()
    }
}