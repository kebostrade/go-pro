package formatter

import (
	"fmt"
	"strings"

	"github.com/DimaJoyti/go-pro/basic/projects/weather-cli/pkg/weather"
)

// Colors for terminal output
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
	Bold   = "\033[1m"
)

// TableFormatter formats weather data as a table
type TableFormatter struct {
	useColors bool
}

// NewTableFormatter creates a new table formatter
func NewTableFormatter(useColors bool) *TableFormatter {
	return &TableFormatter{useColors: useColors}
}

// FormatCurrent formats current weather data
func (f *TableFormatter) FormatCurrent(data *weather.WeatherData, detailed bool) string {
	var sb strings.Builder

	// Header
	sb.WriteString(f.box("╔", "═", "╗", 60))
	sb.WriteString(f.centerText(fmt.Sprintf("Weather in %s, %s",
		data.Location.Name, data.Location.Country), 60))
	sb.WriteString(f.box("╠", "═", "╣", 60))
	sb.WriteString(f.emptyLine(60))

	// Weather icon and condition
	icon := weather.WeatherIcon(strings.ToLower(data.Current.Condition))
	condition := fmt.Sprintf("%s  %s", icon, data.Current.Description)
	sb.WriteString(f.centerText(condition, 60))
	sb.WriteString(f.emptyLine(60))

	// Temperature
	sb.WriteString(f.formatLine("Temperature:",
		fmt.Sprintf("%.1f°F (%.1f°C)", data.Current.Temperature, (data.Current.Temperature-32)*5/9), 60))
	sb.WriteString(f.formatLine("Feels Like:",
		fmt.Sprintf("%.1f°F (%.1f°C)", data.Current.FeelsLike, (data.Current.FeelsLike-32)*5/9), 60))

	// Basic info
	sb.WriteString(f.formatLine("Humidity:", fmt.Sprintf("%d%%", data.Current.Humidity), 60))
	sb.WriteString(f.formatLine("Wind Speed:",
		fmt.Sprintf("%.1f mph (%.1f km/h) %s", data.Current.WindSpeed, data.Current.WindSpeed*1.60934, data.Current.WindDirection), 60))
	sb.WriteString(f.formatLine("Pressure:", fmt.Sprintf("%d hPa", data.Current.Pressure), 60))
	sb.WriteString(f.formatLine("Visibility:", fmt.Sprintf("%d km", data.Current.Visibility), 60))

	if detailed {
		sb.WriteString(f.formatLine("UV Index:", f.formatUVIndex(data.Current.UVIndex), 60))
		sb.WriteString(f.formatLine("Cloud Cover:", fmt.Sprintf("%d%%", data.Current.CloudCover), 60))
		sb.WriteString(f.emptyLine(60))
		sb.WriteString(f.formatLine("Sunrise:", data.Current.Sunrise.Format("03:04 PM"), 60))
		sb.WriteString(f.formatLine("Sunset:", data.Current.Sunset.Format("03:04 PM"), 60))
	}

	sb.WriteString(f.emptyLine(60))
	sb.WriteString(f.formatLine("Last Updated:", data.LastUpdated.Format("2006-01-02 15:04:05"), 60))
	sb.WriteString(f.box("╚", "═", "╝", 60))

	return sb.String()
}

// FormatForecast formats forecast data
func (f *TableFormatter) FormatForecast(data *weather.WeatherData) string {
	var sb strings.Builder

	sb.WriteString(f.box("┌", "─", "┐", 60))
	sb.WriteString(f.centerText(fmt.Sprintf("%d-Day Forecast for %s",
		len(data.Forecast), data.Location.Name), 60))
	sb.WriteString(f.box("├", "─", "┤", 60))
	sb.WriteString(f.emptyLine(60))

	for _, forecast := range data.Forecast {
		icon := weather.WeatherIcon(strings.ToLower(forecast.Condition))
		line := fmt.Sprintf("  %s    %s   %.0f°C / %.0f°C    %s",
			forecast.Date.Format("Mon 02 Jan"),
			icon,
			forecast.TempMax,
			forecast.TempMin,
			forecast.Description,
		)
		sb.WriteString(f.padLine(line, 60))
	}

	sb.WriteString(f.emptyLine(60))
	sb.WriteString(f.box("└", "─", "┘", 60))

	return sb.String()
}

// Helper functions

func (f *TableFormatter) box(left, middle, right string, width int) string {
	return fmt.Sprintf("%s%s%s\n", left, strings.Repeat(middle, width-2), right)
}

func (f *TableFormatter) emptyLine(width int) string {
	return fmt.Sprintf("║%s║\n", strings.Repeat(" ", width-2))
}

func (f *TableFormatter) centerText(text string, width int) string {
	padding := (width - 2 - len(text)) / 2
	return fmt.Sprintf("║%s%s%s║\n",
		strings.Repeat(" ", padding),
		text,
		strings.Repeat(" ", width-2-padding-len(text)))
}

func (f *TableFormatter) formatLine(label, value string, width int) string {
	line := fmt.Sprintf("  %s%s", label, strings.Repeat(" ", 20-len(label)))
	line += value
	return f.padLine(line, width)
}

func (f *TableFormatter) padLine(text string, width int) string {
	if len(text) >= width-2 {
		return fmt.Sprintf("║%s║\n", text[:width-2])
	}
	return fmt.Sprintf("║%s%s║\n", text, strings.Repeat(" ", width-2-len(text)))
}

func (f *TableFormatter) formatUVIndex(uv float64) string {
	level := "Low"
	if uv >= 8 {
		level = "Very High"
	} else if uv >= 6 {
		level = "High"
	} else if uv >= 3 {
		level = "Moderate"
	}
	return fmt.Sprintf("%.0f (%s)", uv, level)
}

func (f *TableFormatter) colorize(text, color string) string {
	if !f.useColors {
		return text
	}
	return color + text + Reset
}
