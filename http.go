package main

import (
	"net/http"
	"io/ioutil"
	"io"
	"encoding/json"
	"log"
	"github.com/gorilla/mux"
)

type Update struct {
	Temperature	float64	`json:"temperature,omitempty"`
	Proportion float64 `json:"proportion,omitempty"`
}

func updateSettings(w http.ResponseWriter, r *http.Request) {
	var update Update
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		if err := r.Body.Close(); err != nil {
    	    panic(err)
    	}

		if err := json.Unmarshal(body, &update); err != nil {
        	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        	w.WriteHeader(422) // unprocessable entity
        	if err := json.NewEncoder(w).Encode(err); err != nil {
        	    panic(err)
        	}
    	}
    	if update.Temperature != 0 {
    		setTargetTemp(update.Temperature)
		}
		if update.Proportion != 0 {
    		log.Printf("updating p constant to %.2f\n", update.Proportion)
			p = update.Proportion
		}
		
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func setpointHandler() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/update", updateSettings)
	log.Fatal(http.ListenAndServe(":2113", router))
}
