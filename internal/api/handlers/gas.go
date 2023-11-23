package handlers

import (
	"encoding/json"
	"ftm-explorer/internal/logger"
	"ftm-explorer/internal/repository"
	"math"
	"net/http"
)

// GasPriceHandler returns the current gas price.
func GasPriceHandler(repo repository.IRepository, log logger.ILogger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		gasPrice, err := repo.SuggestGasPrice()
		if err != nil {
			log.Errorf("failed to get gas price: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		price := math.Round(float64(gasPrice.Int64())/float64(10000000)) / 10.0

		resStruct := struct {
			// Fast is +10% of the average
			Fast float64 `json:"fast"`
			// Fastest is +20% of the average
			Fastest float64 `json:"fastest"`
			SafeLow float64 `json:"safeLow"`
			Average float64 `json:"average"`
		}{
			Fast:    math.Round(price*11) / 10.0,
			Fastest: math.Round(price*12) / 10.0,
			SafeLow: price,
			Average: price,
		}

		res, err := json.Marshal(resStruct)
		if err != nil {
			log.Errorf("failed to marshal gas price: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(res)
	}
}
