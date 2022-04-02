package core

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func init() {

}


// PersistTemp creates a file in /dev/shm
// and saves the setppoint temperature to it. Every time
// the setupoint is changed, this function should be called
// and will persist the value. This is not persistant across
// system reboots.
func PersistTemp(currentTemp float64) {
	persistFile := "/dev/shm/persistFile"
	f, err := os.OpenFile(persistFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("persistence file write failed : %s", err)
	}
	defer f.Close()
	fmt.Fprintf(f, "%f", currentTemp)
	f.Sync()
}

// RecoverTemp is a companion to PersistTemp,
// reading from the file in /dev/shm and setting the setpoint
// at application startup.
func RecoverTemp() (recoverTemp float64) {
	recoverFile := "/dev/shm/persistFile"

	log.Printf("reading setpoint from %s", recoverFile)

	f, err := os.ReadFile(recoverFile)
	if err != nil {
		log.Fatalf("persistence file recovery failed: %s", err)
	}

	recoverTemp, err = strconv.ParseFloat(string(f), 64)
	if err != nil {
		log.Fatalf("persistence conversion failed: %s", err)
	}

	log.Printf("setting setpoint to %f", recoverTemp)
	return recoverTemp

}
