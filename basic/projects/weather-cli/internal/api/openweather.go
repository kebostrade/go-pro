package api

import (
	"context"
	"fmt"
	"time"

	"github.com/DimaJoyti/go-pro/basic/projects/weather-cli/pkg/weather"
)

const openWeatherBaseURL = "https://api.openweathermap.org/data/2.5"

// OpenWeatherClient implements Client for OpenWeatherMap API
type OpenWeatherClient struct {
	*BaseClient
}

// NewOpenWeatherClient creates a new OpenWeatherMap client
func NewOpenWeatherClient(apiKey string) *OpenWeatherClient {
	config := Config{
		APIKey:  apiKey,
		BaseURL: openWeatherBaseURL,
		Timeout: 10 * time.Second,
	}

	return &OpenWeatherClient{
		BaseClient: NewBaseClient(config),
	}
}

// OpenWeatherResponse represents the API response structure
type OpenWeatherResponse struct {
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Sys struct {
		Country string `json:"country"`
		Sunrise int64  `json:"sunrise"`
		Sunset  int64  `json:"sunset"`
	} `json:"sys"`
	Visibility int    `json:"visibility"`
	Name       string `json:"name"`
	Coord      struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"coord"`
}

// GetCurrent retrieves current weather for a city
func (c *OpenWeatherClient) GetCurrent(ctx context.Context, city string, units weather.Units) (*weather.WeatherData, error) {
	url := fmt.Sprintf("%s/weather?q=%s&appid=%s&units=%s",
		c.config.BaseURL, city, c.config.APIKey, units)

	data, err := c.doRequest(ctx, url)
	if err != nil {
		return nil, err
	}

	var resp OpenWeatherResponse
	if err := parseJSON(data, &resp); err != nil {
		return nil, err
	}

	return c.convertToWeatherData(&resp), nil
}

// GetCurrentByCoords retrieves current weather by coordinates
func (c *OpenWeatherClient) GetCurrentByCoords(ctx context.Context, lat, lon float64, units weather.Units) (*weather.WeatherData, error) {
	url := fmt.Sprintf("%s/weather?lat=%f&lon=%f&appid=%s&units=%s",
		c.config.BaseURL, lat, lon, c.config.APIKey, units)

	data, err := c.doRequest(ctx, url)
	if err != nil {
		return nil, err
	}

	var resp OpenWeatherResponse
	if err := parseJSON(data, &resp); err != nil {
		return nil, err
	}

	return c.convertToWeatherData(&resp), nil
}

// GetForecast retrieves weather forecast
func (c *OpenWeatherClient) GetForecast(ctx context.Context, city string, days int, units weather.Units) (*weather.WeatherData, error) {
	// For now, we'll use the current weather endpoint and generate a simple forecast
	// In a real implementation, you would parse the forecast API response
	current, err := c.GetCurrent(ctx, city, units)
	if err != nil {
		return nil, err
	}

	// Add forecast data (simplified for demo)
	current.Forecast = make([]weather.Forecast, days)
	for i := 0; i < days; i++ {
		current.Forecast[i] = weather.Forecast{
			Date:        time.Now().AddDate(0, 0, i+1),
			TempMax:     current.Current.Temperature + float64(i),
			TempMin:     current.Current.Temperature - float64(i+2),
			Condition:   current.Current.Condition,
			Description: current.Current.Description,
		}
	}

	return current, nil
}

// convertToWeatherData converts API response to internal model
func (c *OpenWeatherClient) convertToWeatherData(resp *OpenWeatherResponse) *weather.WeatherData {
	condition := "Unknown"
	description := ""
	if len(resp.Weather) > 0 {
		condition = resp.Weather[0].Main
		description = resp.Weather[0].Description
	}

	return &weather.WeatherData{
		Location: weather.Location{
			Name:      resp.Name,
			Country:   resp.Sys.Country,
			Latitude:  resp.Coord.Lat,
			Longitude: resp.Coord.Lon,
		},
		Current: weather.Current{
			Temperature:   resp.Main.Temp,
			FeelsLike:     resp.Main.FeelsLike,
			Condition:     condition,
			Description:   description,
			Humidity:      resp.Main.Humidity,
			Pressure:      resp.Main.Pressure,
			WindSpeed:     resp.Wind.Speed,
			WindDirection: degToDirection(resp.Wind.Deg),
			Visibility:    resp.Visibility / 1000, // Convert to km
			CloudCover:    resp.Clouds.All,
			Sunrise:       time.Unix(resp.Sys.Sunrise, 0),
			Sunset:        time.Unix(resp.Sys.Sunset, 0),
		},
		LastUpdated: time.Now(),
	}
}

// degToDirection converts wind degree to direction
func degToDirection(deg int) string {
	directions := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	index := int((float64(deg)+22.5)/45.0) % 8
	return directions[index]
}
