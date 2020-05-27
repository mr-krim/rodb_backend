package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func apiResponse(w http.ResponseWriter, r *http.Request) {
	// Set the return Content-Type as JSON like before
	w.Header().Set("Content-Type", "application/json")

	// Change the response depending on the method being requested
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "GET method requested"}`))
	case "POST":
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "POST method requested"}`))
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Can't find method requested"}`))
	}
}

func checkReqEnvVars() {
	checkEnvVars := func(key string) {
		val, ok := os.LookupEnv(key)
		if !ok {
			fmt.Printf("%s is not set\n", key)
			// log.Fatal("Exception: value is not set")
		} else {
			fmt.Printf("%s=%s\n", key, val)
		}
	}
	checkEnvVars("RODB_DEBUG")
	checkEnvVars("RODB_DB_HOST")
	checkEnvVars("RODB_DB_USERNAME")
	checkEnvVars("RODB_DB_PASSWORD")
}

func main() {
	checkReqEnvVars()
	http.HandleFunc("/", apiResponse)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
