package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PolicyHandler struct {
	repo *PolicyRepository
}

func (p PolicyHandler) ListPolicies(w http.ResponseWriter, r *http.Request) {
	var token = r.URL.Query().Get("nextToken")
	fmt.Println("Got here with listpolicies token: ", token)
	var values, nextToken, err = p.repo.listPolicies(token)

	if err == nil {
		err = json.NewEncoder(w).Encode(map[string]any{
			"policies":  values,
			"nextToken": nextToken,
		})
	}

	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}

func (p PolicyHandler) GetPolicies(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "name")
	fmt.Println("Got get request for policy with name:", id)
	fmt.Println(id)
	policy, error := p.repo.getPolicyByName(id)
	if error == nil {
		err := json.NewEncoder(w).Encode(policy)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, error.Error(), http.StatusNotFound)
	}

}

func (p PolicyHandler) CreatePolicy(w http.ResponseWriter, r *http.Request) {
	var policy Policy
	fmt.Println("got here 1")
	var id = uuid.New().String()
	err := json.NewDecoder(r.Body).Decode(&policy)
	policy.ID = id

	fmt.Printf("Decoded policy object %s\n", policy)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to decode request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p.repo.addPolicy(policy)

	err = json.NewEncoder(w).Encode(policy)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

}
func (p PolicyHandler) EnablePolicy(w http.ResponseWriter, r *http.Request) {}

func (p PolicyHandler) DeletePolicy(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "name")
	p.repo.deletePolicy(id)

}

func (p PolicyHandler) UpdatePolicy(w http.ResponseWriter, r *http.Request) {
	var policy Policy
	fmt.Printf("update got here %s\n", chi.URLParam(r, "name"))
	err := json.NewDecoder(r.Body).Decode(&policy)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to decode request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("update got here 2")

	p.repo.updatePolicy(policy)

	fmt.Println("got here 3")

	err = json.NewEncoder(w).Encode(policy)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

}
