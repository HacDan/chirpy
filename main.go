package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type apiConfig struct {
	fileserverHits int
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}
type CleanedChirp struct {
	CleanedBody string `json:"cleaned_body"`
}

type RespError struct {
	Error string `json:"error"`
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
	if strings.Contains(strings.ToLower(chirp.Body), "kerfuffle") || strings.Contains(strings.ToLower(chirp.Body), "sharbert") || strings.Contains(strings.ToLower(chirp.Body), "sharbert") {
		splitChirp := strings.Split(chirp.Body, " ")

		for i, word := range splitChirp {
			if strings.ToLower(word) == "kerfuffle" {
				splitChirp[i] = "****"
			}
			if strings.ToLower(word) == "sharbert" {
				splitChirp[i] = "****"
			}
			if strings.ToLower(word) == "fornax" {
				splitChirp[i] = "****"
			}
		}

		cleanedChirp := CleanedChirp{
			CleanedBody: strings.Join(splitChirp, " "),
		}

		dat, err := json.Marshal(cleanedChirp)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(dat)
		return
	}

	cleanedChirp := CleanedChirp{
		CleanedBody: chirp.Body,
	}
	dat, err := json.Marshal(cleanedChirp)
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
