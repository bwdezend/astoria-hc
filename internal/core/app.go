package core

import (
	"log"
	"time"

	"github.com/brutella/hc/accessory"
	"github.com/bwdezend/astoria-hc/internal/telemetry"
	"github.com/vemo-france/max31865"
)

var p = 3.0

// GetCurrentTemp doc
func GetCurrentTemp(acc accessory.Thermostat) {
	if err := max31865.Init(); err != nil {
		log.Fatalf("initialization failed : %s", err)
	}
	sensor := max31865.Create("8", "9", "10", "11")
	var boilerTemperature float64
	for {
		boilerTemperature = float64(sensor.ReadTemperature(100, 430))
		acc.Thermostat.CurrentTemperature.SetValue(boilerTemperature)
		telemetry.CurrentTemperature.Set(boilerTemperature)
		time.Sleep(500 * time.Millisecond)
	}

}

// SetTargetTemp adjusts the setpoint for the PID loop and updates the homekit interfaces
func SetTargetTemp(acc accessory.Thermostat, setTemp float64) {
	if setTemp > 124.0 {
		setTemp = 124.0
	}

	log.Printf("setting setpoint to %.2f", setTemp)
	acc.Thermostat.TargetTemperature.SetValue(setTemp)
	telemetry.SetpointTemperature.Set(setTemp)
}

// HeaterWindow doc Take two inputs - the duration of the cycle and the proportion of the cycle
// and turn the heating element on for that percentage of the cycle. If the
// windowSize is 15.0 and the enabledTime is 0.7, this turns the heating element
// on for 10.5 seconds and off for 4.5 seconds before returning
func HeaterWindow(windowSize float64, enabledTime float64, gpioEnabled bool) {
	var disabledTime float64 = 0
	enabledTime = windowSize * enabledTime * 1000
	disabledTime = windowSize*1000 - enabledTime

	if enabledTime > 1000.0 {
		enabledTime = 1000.0
	}
	if enabledTime > 0 {
		if gpioEnabled {
			HeaterControl(true)
		}
		time.Sleep(time.Duration(enabledTime) * time.Millisecond)
		telemetry.SecondsActive.Add(enabledTime / 1000)
	}
	if gpioEnabled {
		HeaterControl(false)
	}
	time.Sleep(time.Duration(disabledTime) * time.Millisecond)
}

// TemperatureProportional is a dead simple proportional control loop. Take the difference in setpoint
// and current temp, multiply by the gain, and use the result to control
// the duty cycle on the boiler, represented as a float between 0.0 and 1.0
func TemperatureProportional(acc accessory.Thermostat, gpioEnabled bool) {

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

		HeaterWindow(1.0, error, gpioEnabled)

	}
}

// TemperatureErrorDetection is a error handling function to turn the SSR off if the boiler doesn't
// appear to be warming as a result of input. The most common reason for this
// is that the power switch on the machine is turned off, but the timer still
// fired. Less common reasons would be a heater element malfunction, a boiler
// rupture, or a failure of the temperature sensor itself.
func TemperatureErrorDetection(acc accessory.Thermostat) {

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
					SetTargetTemp(acc, 0.0)
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

// Re-evaluate if we want to be able to set the temp from a REST endpoint
// type updateVars struct {
// 	Temperature float64 `json:"temperature,omitempty"`
// 	Proportion  float64 `json:"proportion,omitempty"`
// }
//
// // UpdateSettings doc
// func UpdateSettings(w http.ResponseWriter, r *http.Request) {
// 	var update updateVars
// 	if r.Method == "POST" {
// 		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
// 		if err != nil {
// 			http.Error(w, "Error reading request body",
// 				http.StatusInternalServerError)
// 		}
// 		if err := r.Body.Close(); err != nil {
// 			panic(err)
// 		}
//
// 		if err := json.Unmarshal(body, &update); err != nil {
// 			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 			w.WriteHeader(422) // unprocessable entity
// 			if err := json.NewEncoder(w).Encode(err); err != nil {
// 				panic(err)
// 			}
// 		}
// 		if update.Temperature != 0 {
// 			SetTargetTemp(update.Temperature)
// 		}
// 		if update.Proportion != 0 {
// 			log.Printf("updating p constant to %.2f\n", update.Proportion)
// 			p = update.Proportion
// 		}
//
// 	} else {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 	}
// }
//
// // SetpointHandler handles
// func SetpointHandler() {
// 	router := mux.NewRouter().StrictSlash(true)
// 	router.HandleFunc("/update", UpdateSettings)
// 	log.Fatal(http.ListenAndServe(":2113", router))
// }
