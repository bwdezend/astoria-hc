package telemetry

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// SecondsActive is the number of seconds, since program start
	//  that the heating elemnent has been active.
	SecondsActive = promauto.NewCounter(prometheus.CounterOpts{
		Name: "heater_element_active_seconds",
		Help: "The total number of seconds the heating element has been active",
	})

	// RelayActivations counts the number of times the relay has activated
	RelayActivations = promauto.NewCounter(prometheus.CounterOpts{
		Name: "heater_element_activations_total",
		Help: "The number of times the relay has been activated",
	})

	// CurrentTemperature is the current measured temperature in the boiler
	CurrentTemperature = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "boiler_water_temperature_celsius",
		Help: "The current temperature of the water in the boiler",
	})

	// SetpointTemperature is the desired temperature
	SetpointTemperature = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "setpoint_temperature_celsius",
		Help: "The current setpoint temperature",
	})
)

func exitProgram(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}

// PrometheusMetrics handles the prometheus scrape endpoint
func PrometheusMetrics(prometheusPort int) {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/exit", exitProgram)
	http.ListenAndServe(fmt.Sprintf(":%d", prometheusPort), nil)
}
