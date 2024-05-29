package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

type Chirp struct {
	Body string `json:"body"`
}

type RespError struct {
	Error string `json:"error"`
}

type RespValid struct {
	Valid bool `json:"valid"`
}

func main() {
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	apiCfg := apiConfig{}

	mux.HandleFunc("GET /api/healthz", HandleHealthz)
	mux.HandleFunc("GET /api/reset", apiCfg.HandleReset)
	mux.HandleFunc("POST /api/validate_chirp", HandleValidateChirp)
	mux.HandleFunc("GET /admin/metrics", apiCfg.HandleMetrics)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func HandleHealthz(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func HandleValidateChirp(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	chirp := Chirp{}

	err := decoder.Decode(&chirp)
	if err != nil {
		respError := RespError{
			Error: "Something went wrong",
		}
		dat, err := json.Marshal(respError)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(dat)
		return
	}

	if len(chirp.Body) > 140 {
		respError := RespError{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(respError)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	respValid := RespValid{
		Valid: true,
	}

	dat, err := json.Marshal(respValid)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)

}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) HandleMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	metricsResponse := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits)
	w.Write([]byte(metricsResponse))
}

func (cfg *apiConfig) HandleReset(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits = 0
}
