// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package domain defines analytics domain models.
package domain

import "time"

// StudentMetrics represents analytics data for a student.
type StudentMetrics struct {
	UserID                string                 `json:"user_id"`
	LearningVelocity      *LearningVelocityData `json:"learning_velocity"`
	ErrorPatterns         []ErrorPattern         `json:"error_patterns"`
	TimeDistribution      []TopicTimeSpent       `json:"time_distribution"`
	PeerComparison        *PeerComparisonData    `json:"peer_comparison"`
	WeakAreas             []WeakArea             `json:"weak_areas"`
	StudyRecommendations  []string               `json:"study_recommendations"`
	GeneratedAt           time.Time              `json:"generated_at"`
}

// LearningVelocityData represents lessons completed over time.
type LearningVelocityData struct {
	LessonsPerWeek     float64     `json:"lessons_per_week"`
	TwelveWeekTrend    []WeekData  `json:"twelve_week_trend"`
	AverageTimePerLesson time.Duration `json:"average_time_per_lesson"`
}

// WeekData represents weekly progress data.
type WeekData struct {
	Week         int       `json:"week"`
	Lessons      int       `json:"lessons"`
	WeekStart    time.Time `json:"week_start"`
	WeekEnd      time.Time `json:"week_end"`
}

// ErrorPattern represents a common mistake pattern.
type ErrorPattern struct {
	Topic       string  `json:"topic"`
	ErrorCount  int     `json:"error_count"`
	Frequency   float64 `json:"frequency"` // Percentage of attempts
	LastSeen    time.Time `json:"last_seen"`
}

// TopicTimeSpent represents time spent on a topic.
type TopicTimeSpent struct {
	Topic        string        `json:"topic"`
	TimeSpent    time.Duration `json:"time_spent"`
	Percentage   float64       `json:"percentage"`
	LessonCount  int           `json:"lesson_count"`
}

// PeerComparisonData represents comparison with class performance.
type PeerComparisonData struct {
	UserPercentile    float64   `json:"user_percentile"`
	ClassAverage      float64   `json:"class_average"`
	ClassDistribution []float64 `json:"class_distribution"` // Percentiles
	UserScore         float64   `json:"user_score"`
}

// WeakArea represents a topic needing review.
type WeakArea struct {
	Topic          string  `json:"topic"`
	AverageScore   float64 `json:"average_score"`
	ClassAverage   float64 `json:"class_average"`
	Severity       string  `json:"severity"` // "critical", "warning", "info"
	RecommendedLessons []string `json:"recommended_lessons"`
}

// InstructorMetrics represents analytics data for an instructor.
type InstructorMetrics struct {
	ClassID              string                `json:"class_id"`
	ClassOverview        *ClassOverviewData    `json:"class_overview"`
	EngagementHeatmap    []EngagementDataPoint `json:"engagement_heatmap"`
	CompletionFunnel     *CompletionFunnelData `json:"completion_funnel"`
	TimeToCompletion     []TimeDistribution    `json:"time_to_completion"`
	StrugglingStudents   []StrugglingStudent    `json:"struggling_students"`
	ContentPerformance   []ContentMetric       `json:"content_performance"`
	GradeDistribution    []GradeBucket         `json:"grade_distribution"`
	GeneratedAt          time.Time             `json:"generated_at"`
}

// ClassOverviewData represents high-level class metrics.
type ClassOverviewData struct {
	TotalStudents      int     `json:"total_students"`
	ActiveThisWeek     int     `json:"active_this_week"`
	AverageCompletion  float64 `json:"average_completion"`
	AverageScore       float64 `json:"average_score"`
}

// EngagementDataPoint represents activity at a specific time.
type EngagementDataPoint struct {
	DayOfWeek   int     `json:"day_of_week"` // 0-6 (Sunday-Saturday)
	HourOfDay   int     `json:"hour_of_day"` // 0-23
	ActivityCount int  `json:"activity_count"`
}

// CompletionFunnelData represents progression through assessments.
type CompletionFunnelData struct {
	AssessmentID    string `json:"assessment_id"`
	AssessmentTitle string `json:"assessment_title"`
	Started         int    `json:"started"`
	Submitted       int    `json:"submitted"`
	Passed          int    `json:"passed"`
}

// TimeDistribution represents completion time distribution.
type TimeDistribution struct {
	BucketLowerBound time.Duration `json:"bucket_lower_bound"`
	BucketUpperBound time.Duration `json:"bucket_upper_bound"`
	StudentCount     int            `json:"student_count"`
	Percentage       float64        `json:"percentage"`
}

// StrugglingStudent represents a student needing intervention.
type StrugglingStudent struct {
	UserID        string    `json:"user_id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	AverageScore  float64   `json:"average_score"`
	ClassAverage  float64   `json:"class_average"`
	LastActive    time.Time `json:"last_active"`
	InactiveDays  int       `json:"inactive_days"`
	Reason        []string  `json:"reason"` // ["low_score", "inactive"]
}

// ContentMetric represents performance metrics for content.
type ContentMetric struct {
	ContentType   string  `json:"content_type"` // "lesson" or "assessment"
	ContentID     string  `json:"content_id"`
	ContentTitle  string  `json:"content_title"`
	PassRate      float64 `json:"pass_rate"`
	AverageScore  float64 `json:"average_score"`
	AverageTimeSpent time.Duration `json:"average_time_spent"`
	Engagement    float64 `json:"engagement"` // Percentage of students who interacted
}

// GradeBucket represents score distribution.
type GradeBucket struct {
	RangeLabel string  `json:"range_label"` // "0-20", "21-40", etc.
	MinScore   float64 `json:"min_score"`
	MaxScore   float64 `json:"max_score"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// AdminMetrics represents platform-wide analytics for admins.
type AdminMetrics struct {
	PlatformMetrics   *PlatformMetricsData `json:"platform_metrics"`
	UserGrowth        *UserGrowthData      `json:"user_growth"`
	ContentPopularity []ContentMetric      `json:"content_popularity"`
	SystemHealth      *SystemHealthData    `json:"system_health"`
	StorageUsage      *StorageUsageData    `json:"storage_usage"`
	SandboxUsage      *SandboxUsageData    `json:"sandbox_usage"`
	GeneratedAt       time.Time            `json:"generated_at"`
}

// PlatformMetricsData represents overall platform statistics.
type PlatformMetricsData struct {
	TotalUsers     int64   `json:"total_users"`
	DAU            int64   `json:"dau"` // Daily Active Users
	WAU            int64   `json:"wau"` // Weekly Active Users
	MAU            int64   `json:"mau"` // Monthly Active Users
	ContentCount   int64   `json:"content_count"`
	NewSignupsThisWeek int64 `json:"new_signups_this_week"`
}

// UserGrowthData represents user growth over time.
type UserGrowthData struct {
	SignupsTrend    []GrowthDataPoint `json:"signups_trend"`    // 12-month trend
	ChurnRate       float64           `json:"churn_rate"`        // Percentage
	RetentionRate   float64           `json:"retention_rate"`    // Percentage
	CohortAnalysis  []CohortData      `json:"cohort_analysis"`   // Retention by signup month
}

// GrowthDataPoint represents growth at a point in time.
type GrowthDataPoint struct {
	Month        time.Time `json:"month"`
	UserCount    int64     `json:"user_count"`
	NewSignups   int64     `json:"new_signups"`
}

// CohortData represents user retention by signup cohort.
type CohortData struct {
	CohortMonth      time.Time `json:"cohort_month"`
	InitialUsers     int       `json:"initial_users"`
	RetentionWeek1   float64   `json:"retention_week_1"`
	RetentionWeek4   float64   `json:"retention_week_4"`
	RetentionWeek12  float64   `json:"retention_week_12"`
}

// SystemHealthData represents system performance metrics.
type SystemHealthData struct {
	APIResponseTimeP50  time.Duration `json:"api_response_time_p50"`
	APIResponseTimeP95  time.Duration `json:"api_response_time_p95"`
	APIResponseTimeP99  time.Duration `json:"api_response_time_p99"`
	ErrorRate          float64       `json:"error_rate"`      // Percentage
	DBQueryTimeAvg      time.Duration `json:"db_query_time_avg"`
	CacheHitRate       float64       `json:"cache_hit_rate"`   // Percentage
	Uptime             float64       `json:"uptime"`           // Percentage
}

// StorageUsageData represents storage consumption.
type StorageUsageData struct {
	TotalStorageUsedGB    float64            `json:"total_storage_used_gb"`
	VideosStorageGB       float64            `json:"videos_storage_gb"`
	DownloadsStorageGB    float64            `json:"downloads_storage_gb"`
	SubmissionsStorageGB  float64            `json:"submissions_storage_gb"`
	ProjectedGrowth       []GrowthProjection `json:"projected_growth"`
}

// GrowthProjection projects storage usage over time.
type GrowthProjection struct {
	Month       time.Time `json:"month"`
	ProjectedGB float64   `json:"projected_gb"`
}

// SandboxUsageData represents code execution sandbox metrics.
type SandboxUsageData struct {
	TotalExecutions       int64         `json:"total_executions"`
	SuccessRate           float64       `json:"success_rate"`
	AverageExecutionTime  time.Duration `json:"average_execution_time"`
	AverageMemoryUsageMB  float64       `json:"average_memory_usage_mb"`
	ErrorDistribution     []ErrorMetric `json:"error_distribution"`
}

// ErrorMetric represents types of sandbox errors.
type ErrorMetric struct {
	ErrorType   string  `json:"error_type"` // "compile_error", "runtime_error", "timeout"
	Count       int     `json:"count"`
	Percentage  float64 `json:"percentage"`
}

// AnalyticsQueryParams represents query parameters for analytics requests.
type AnalyticsQueryParams struct {
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	UserID      string     `json:"user_id,omitempty"`
	ClassID     string     `json:"class_id,omitempty"`
	LessonID    string     `json:"lesson_id,omitempty"`
	ExportFormat string    `json:"export_format,omitempty"` // "json" or "csv"
}

// ExportRequest represents a request to export analytics data.
type ExportRequest struct {
	MetricsType   string            `json:"metrics_type"`   // "student", "instructor", "admin"
	QueryParams   AnalyticsQueryParams `json:"query_params"`
	Format        string            `json:"format"`         // "json" or "csv"
}
