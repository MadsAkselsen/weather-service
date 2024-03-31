package main

import (
	"log"
	"net/http"
	"weather-service/config"
	"weather-service/handlers"
	"weather-service/store"

	"github.com/go-redis/redis/v8"
)

func main() {
	// Load config
	cfg := config.NewConfig()

	// Initialize Redis client
	redisClient := store.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: "", // No password set
		DB:       0,  // Use default DB
	})
	cacheStore := store.NewStore(redisClient)

	// Initialize handlers with the cache store and config
	handler := handlers.NewHandler(cacheStore, cfg)

	// Routes
	http.Handle("/weather", corsMiddleware(http.HandlerFunc(handler.GetWeatherData)))
	http.Handle("/weatherByCoords", corsMiddleware(http.HandlerFunc(handler.GetWeatherByCoords)))

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// corsMiddleware adds CORS headers to responses
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                   // Allow any domain, adjust if necessary for production
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // Adjust the methods as per your requirements
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")       // Add other headers as needed

		// Check if the request is for the OPTIONS method (pre-flight request)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
