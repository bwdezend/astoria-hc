package core

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func init() {

}

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

func RecoverTemp() (recoverTemp float64) {
	//recoverTemp = 0.0
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
