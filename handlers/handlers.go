package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"weather-service/config"
	"weather-service/store"
)

type Handler struct {
	store  *store.Store
	config *config.Config
}

// Adjust NewHandler to accept the new Config parameter
func NewHandler(store *store.Store, cfg *config.Config) *Handler {
	return &Handler{
		store:  store,
		config: cfg,
	}
}

// WeatherData is a simplified structure to hold weather data.
// Adjust it according to the actual data you're working with.
type WeatherData struct {
	CityName    string `json:"cityName"`
	Temperature string `json:"temperature"`
	// Include other fields as needed
}

func (h *Handler) GetWeatherData(w http.ResponseWriter, r *http.Request) {
	cityName := r.URL.Query().Get("city")
	if cityName == "" {
		http.Error(w, "City name is required", http.StatusBadRequest)
		return
	}

	cacheKey := "weather_" + cityName
	cachedData, err := h.store.Get(cacheKey)
	fmt.Println("Cached data:", cachedData)
	if err != nil {
		// Fetch current weather
		weatherApiUrl := h.config.OpenWeatherMapURL + "/weather?q=" + cityName + "&appid=" + h.config.WeatherAPIKey + "&units=metric"
		weatherData, err := fetchDataFromAPI(weatherApiUrl, map[string]string{"Content-Type": "application/json"})
		if err != nil {
			fmt.Println("Failed to fetch weather data:", err)
			http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
			return
		}

		// Fetch forecast
		forecastApiUrl := h.config.OpenWeatherMapURL + "/forecast?q=" + cityName + "&appid=" + h.config.WeatherAPIKey + "&units=metric"
		forecastData, err := fetchDataFromAPI(forecastApiUrl, map[string]string{"Content-Type": "application/json"})
		if err != nil {
			fmt.Println("Failed to fetch forecast data:", err)
			http.Error(w, "Failed to fetch forecast data", http.StatusInternalServerError)
			return
		}

		// Combine weather and forecast data
		combinedData := fmt.Sprintf(`{"weatherData": %s, "forecastData": %s}`, string(weatherData), string(forecastData))

		// Cache the combined data
		h.store.Set(cacheKey, combinedData, 24*time.Hour)

		// Respond with combined data
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(combinedData))
	} else {
		// Send cached data if available
		fmt.Println("Sending cached data")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cachedData))
	}
}

func (h *Handler) GetWeatherByCoords(w http.ResponseWriter, r *http.Request) {
	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("lon")
	if lat == "" || lon == "" {
		http.Error(w, "Latitude and longitude are required", http.StatusBadRequest)
		return
	}

	cacheKey := "weather_coords_" + lat + "_" + lon
	cachedData, err := h.store.Get(cacheKey)
	fmt.Println("Cached data:", cachedData)
	if err != nil {
		// Fetch current weather data
		weatherApiUrl := h.config.OpenWeatherMapURL + "/weather?lat=" + lat + "&lon=" + lon + "&appid=" + h.config.WeatherAPIKey + "&units=metric"
		weatherData, err := fetchDataFromAPI(weatherApiUrl, map[string]string{"Content-Type": "application/json"})
		if err != nil {
			fmt.Println("Failed to fetch weather data:", err)
			http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
			return
		}

		// Fetch forecast data
		forecastApiUrl := h.config.OpenWeatherMapURL + "/forecast?lat=" + lat + "&lon=" + lon + "&appid=" + h.config.WeatherAPIKey + "&units=metric"
		forecastData, err := fetchDataFromAPI(forecastApiUrl, map[string]string{"Content-Type": "application/json"})
		if err != nil {
			fmt.Println("Failed to fetch forecast data:", err)
			http.Error(w, "Failed to fetch forecast data", http.StatusInternalServerError)
			return
		}

		// Combine weather and forecast data
		combinedData := fmt.Sprintf(`{"weatherData": %s, "forecastData": %s}`, string(weatherData), string(forecastData))

		// Cache the combined data
		h.store.Set(cacheKey, combinedData, 24*time.Hour)

		// Respond with combined data
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(combinedData))
	} else {
		// Send cached data if available
		fmt.Println("Sending cached data")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cachedData))
	}
}

// fetchDataFromAPI is a utility to fetch data from the given URL with provided headers.
func fetchDataFromAPI(url string, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Set headers
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
