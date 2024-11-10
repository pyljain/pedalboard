package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func NewFromFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.NewDecoder(f).Decode(&config)
	if err != nil {
		return nil, err
	}

	// Check if probabilities add up to 100 per route
	for _, route := range config.Routes {
		sumOfProbabilities := 0
		for _, sampleResponse := range route.Responses {
			sumOfProbabilities += sampleResponse.Probability
		}

		if sumOfProbabilities == 100 {
			continue
		} else {
			return nil, fmt.Errorf("For each route the sum of probabilities of sample response distribution should be 100. This is not the case for %s", route.Path)
		}
	}

	return &config, nil
}
