package main

import (
	"net/http"
	"os"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)


var (
	secondsActive = promauto.NewCounter(prometheus.CounterOpts{
		Name: "heater_element_active_seconds",
		Help: "The total number of seconds the heating element has been active",
	})
	
	relayActivations = promauto.NewCounter(prometheus.CounterOpts{
		Name: "heater_element_activations_total",
		Help: "The number of times the HS100 relay has been activated",
	})

	currentTemperature = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "boiler_water_temperature_celsius",
		Help: "The current temperature of the water in the boiler",
	})

	setpointTemperature = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "setpoint_temperature_celsius",
		Help: "The current setpoint temperature",
	})
)

func exitProgram(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}

func prometheusMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/exit", exitProgram)
	http.ListenAndServe(":2112", nil)
}

