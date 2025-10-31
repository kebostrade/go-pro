package exercises

// Exercise 2: Constants and iota Practice
// Complete the constant declarations and functions below

// Mathematical constants
const (
	Pi          = 3.14159
	E           = 2.71828
	GoldenRatio = 1.61803
)

// HTTP status codes
const (
	StatusOK        = 200
	StatusCreated   = 201
	StatusAccepted  = 202
	StatusNoContent = 204
)

// Log levels using iota
const (
	LogDebug   = iota // 0
	LogInfo           // 1
	LogWarning        // 2
	LogError          // 3
	LogFatal          // 4
)

// File permissions using bit flags
const (
	PermissionRead    = 1 << iota // 1 (binary: 001)
	PermissionWrite               // 2 (binary: 010)
	PermissionExecute             // 4 (binary: 100)
)

// Weekday constants starting from 1
const (
	Monday = iota + 1 // 1
	Tuesday           // 2
	Wednesday         // 3
	Thursday          // 4
	Friday            // 5
	Saturday          // 6
	Sunday            // 7
)

// GetHTTPStatusMessage returns a message for the given HTTP status code
func GetHTTPStatusMessage(statusCode int) string {
	switch statusCode {
	case StatusOK:
		return "OK"
	case StatusCreated:
		return "Created"
	case StatusAccepted:
		return "Accepted"
	case StatusNoContent:
		return "No Content"
	default:
		return "Unknown Status"
	}
}

// GetLogLevelName returns the name of the log level
func GetLogLevelName(level int) string {
	switch level {
	case LogDebug:
		return "DEBUG"
	case LogInfo:
		return "INFO"
	case LogWarning:
		return "WARNING"
	case LogError:
		return "ERROR"
	case LogFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// HasPermission checks if the given permissions include the specified permission
func HasPermission(permissions, permission int) bool {
	return permissions&permission != 0
}

// AddPermission adds a permission to the existing permissions
func AddPermission(permissions, permission int) int {
	return permissions | permission
}

// RemovePermission removes a permission from the existing permissions
func RemovePermission(permissions, permission int) int {
	return permissions &^ permission
}

// GetWeekdayName returns the name of the weekday
func GetWeekdayName(day int) string {
	switch day {
	case Monday:
		return "Monday"
	case Tuesday:
		return "Tuesday"
	case Wednesday:
		return "Wednesday"
	case Thursday:
		return "Thursday"
	case Friday:
		return "Friday"
	case Saturday:
		return "Saturday"
	case Sunday:
		return "Sunday"
	default:
		return "Invalid Day"
	}
}

// IsWeekend checks if the given day is a weekend day
func IsWeekend(day int) bool {
	return day == Saturday || day == Sunday
}

// CalculateCircleProperties calculates area and circumference of a circle
func CalculateCircleProperties(radius float64) (float64, float64) {
	area := Pi * radius * radius
	circumference := 2 * Pi * radius
	return area, circumference
}

// FormatPermissions returns a string representation of permissions
func FormatPermissions(permissions int) string {
	result := ""

	if HasPermission(permissions, PermissionRead) {
		result += "r"
	} else {
		result += "-"
	}

	if HasPermission(permissions, PermissionWrite) {
		result += "w"
	} else {
		result += "-"
	}

	if HasPermission(permissions, PermissionExecute) {
		result += "x"
	} else {
		result += "-"
	}

	return result
}
