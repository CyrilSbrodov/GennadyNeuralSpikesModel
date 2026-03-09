package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"gennady-neural-spikes-model/internal/engine"
)

type app struct {
	mu sync.Mutex
	e  *engine.Engine
}

type stepRequest struct {
	Inputs []engine.InputSignal `json:"inputs"`
	Steps  int                  `json:"steps"`
}

type stepResponse struct {
	Tick      int     `json:"tick"`
	FiredEach [][]int `json:"firedEach"`
}

type saveRequest struct {
	Path string `json:"path"`
}

func main() {
	neuronCount := envInt("NEURONS", 1000)
	synapsesPerNeuron := envInt("SYNAPSES_PER_NEURON", 8)
	randomSyn := make([]engine.SynapseDef, 0, neuronCount*synapsesPerNeuron)
	for src := 0; src < neuronCount; src++ {
		for i := 0; i < synapsesPerNeuron; i++ {
			target := (src*31 + i*17 + 13) % neuronCount
			randomSyn = append(randomSyn, engine.SynapseDef{
				Source:  src,
				Target:  target,
				Weight:  10 + float64((src+i)%5),
				Delay:   1 + (i % 4),
				Plastic: true,
			})
		}
	}

	cfg := engine.Config{LeakFactor: 0.97, MinWeight: 0.1, MaxWeight: 50, LearnRate: 0.2, ForgetRate: 0.01, LearnWindowTicks: 4, ForgetAfterInactive: 40}
	e, err := engine.NewEngine(neuronCount, randomSyn, cfg)
	if err != nil {
		log.Fatal(err)
	}

	a := &app{e: e}
	http.HandleFunc("/step", a.handleStep)
	http.HandleFunc("/state", a.handleState)
	http.HandleFunc("/save", a.handleSave)
	http.HandleFunc("/load", a.handleLoad)

	addr := ":8080"
	log.Printf("spiking service listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (a *app) handleStep(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req stepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Steps <= 0 {
		req.Steps = 1
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	firedEach := make([][]int, 0, req.Steps)
	for i := 0; i < req.Steps; i++ {
		in := req.Inputs
		if i > 0 {
			in = nil
		}
		firedEach = append(firedEach, a.e.Step(in))
	}
	writeJSON(w, stepResponse{Tick: a.e.Tick, FiredEach: firedEach})
}

func (a *app) handleState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	a.mu.Lock()
	s := a.e.Snapshot()
	a.mu.Unlock()
	writeJSON(w, s)
}

func (a *app) handleSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req saveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	a.mu.Lock()
	err := a.e.Save(req.Path)
	a.mu.Unlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"status": "ok"})
}

func (a *app) handleLoad(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req saveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	e, err := engine.Load(req.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	a.mu.Lock()
	a.e = e
	a.mu.Unlock()
	writeJSON(w, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func envInt(name string, fallback int) int {
	if v := os.Getenv(name); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
