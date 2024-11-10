package server

import (
	"fmt"
	"math/rand"
	"net/http"
	"pedalboard/internal/config"
	"time"
)

type server struct {
	cfg *config.Config
}

func New(cfg *config.Config) *server {
	return &server{
		cfg: cfg,
	}
}

func (s *server) Start() error {

	mux := http.NewServeMux()
	var responseProbabilitiesTrackerByRoute [][]int

	for i, route := range s.cfg.Routes {
		responseProbabilitiesTrackerByRoute = append(responseProbabilitiesTrackerByRoute, []int{})
		mux.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
			// Authentication
			if s.cfg.Authentication.APIKey != "" {
				apiKey := r.Header.Get("Authorization")
				if len(apiKey) < 7 || apiKey[:7] != "Bearer " {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				apiKey = apiKey[7:]
				if apiKey != s.cfg.Authentication.APIKey {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}

			// Add latency
			latency := rand.Intn(route.Latency.Max-route.Latency.Min+1) + route.Latency.Min
			time.Sleep(time.Duration(latency) * time.Millisecond)

			// Choose a response based on probability & hydrate the slice if not hydrated with sample response indices
			if len(responseProbabilitiesTrackerByRoute[i]) == 0 {
				responseProbabilitiesTrackerByRoute[i] = s.hydrateTrackerWithResponses(i)
			}

			// Select a sample response based on probabilities of the expected response codes to be returned
			var chosenResponse int
			chosenResponse, responseProbabilitiesTrackerByRoute[i] = responseProbabilitiesTrackerByRoute[i][0], responseProbabilitiesTrackerByRoute[i][1:]

			// Write the response
			w.WriteHeader(route.Responses[chosenResponse].Status)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(route.Responses[chosenResponse].Body))
		})
	}

	svr := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.Port),
		Handler: mux,
	}

	return svr.ListenAndServe()
}

func (s *server) hydrateTrackerWithResponses(routeIndex int) []int {
	var responseProbabilitiesTrackerForRoute []int

	for j, sr := range s.cfg.Routes[routeIndex].Responses {
		for k := 0; k < sr.Probability; k++ {
			responseProbabilitiesTrackerForRoute = append(responseProbabilitiesTrackerForRoute, j)
		}
	}

	rand.Shuffle(len(responseProbabilitiesTrackerForRoute), func(i, j int) {
		responseProbabilitiesTrackerForRoute[i], responseProbabilitiesTrackerForRoute[j] = responseProbabilitiesTrackerForRoute[j], responseProbabilitiesTrackerForRoute[i]
	})

	return responseProbabilitiesTrackerForRoute
}
