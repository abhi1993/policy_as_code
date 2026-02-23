package main

import (
	"io"
	"log"
	"net/http"
)

type AuthHandler struct {
	policyServiceURI string
}

func (handler AuthHandler) PerformAuthorization(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(handler.policyServiceURI + "/policies")

	if err != nil {
		log.Fatal("request failed:", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("failed to read response:", err)
	} else {
		log.Printf("response: %s", body)
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}

}
