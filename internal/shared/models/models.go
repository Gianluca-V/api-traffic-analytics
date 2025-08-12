package models

import (
	"time"
)

// HealthCheckResponse representa la respuesta del endpoint de health check
type HealthCheckResponse struct {
    Status    string            `json:"status"`
    Timestamp time.Time         `json:"timestamp"`
    Services  map[string]string `json:"services,omitempty"`
    Version   string            `json:"version"`
}

// ErrorResponse representa una respuesta de error estandarizada
type ErrorResponse struct {
    Error   string `json:"error"`
    Code    int    `json:"code,omitempty"`
    Message string `json:"message"`
    Details interface{} `json:"details,omitempty"`
}

// Helper function para crear timestamps consistentes
func TimeNow() time.Time {
    return time.Now().UTC()
}

// SuccessResponse representa una respuesta exitosa estandarizada
type SuccessResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
}

// =====================================================
// LOCATIONS
// =====================================================
type Location struct {
	ID          string    `gorm:"primaryKey;size:50" json:"id" validate:"required,max=50"`
	Name        string    `gorm:"size:255;not null" json:"name" validate:"required,max=255"`
	Description string    `json:"description,omitempty"`
	Latitude    float64   `gorm:"type:decimal(10,8);not null" json:"latitude" validate:"required"`
	Longitude   float64   `gorm:"type:decimal(11,8);not null" json:"longitude" validate:"required"`
	Address     string    `json:"address,omitempty"`
	City        string    `gorm:"size:100" json:"city,omitempty" validate:"max=100"`
	Country     string    `gorm:"size:100;default:Argentina" json:"country,omitempty" validate:"max=100"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// =====================================================
// TRAFFIC DATA
// =====================================================
type TrafficData struct {
	ID              int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID            string    `gorm:"type:uuid;default:uuid_generate_v4()" json:"uuid"`
	Timestamp       time.Time `gorm:"autoCreateTime" json:"timestamp"`
	LocationID      string    `gorm:"size:50;not null" json:"location_id" validate:"required"`
	VehicleCount    int       `gorm:"not null" json:"vehicle_count" validate:"required,min=0"`
	AverageSpeed    float64   `gorm:"type:decimal(5,2);not null" json:"average_speed" validate:"required,min=0"`
	CongestionLevel string    `gorm:"size:20;not null" json:"congestion_level" validate:"required,oneof=low medium high severe"`
	MaxSpeed        *float64  `gorm:"type:decimal(5,2)" json:"max_speed,omitempty" validate:"omitempty,min=0"`
	MinSpeed        *float64  `gorm:"type:decimal(5,2)" json:"min_speed,omitempty" validate:"omitempty,min=0"`
	Occupancy       *float64  `gorm:"type:decimal(5,2)" json:"occupancy,omitempty" validate:"omitempty,min=0,max=100"`
	QueueLength     float64   `gorm:"type:decimal(8,2);default:0" json:"queue_length"`
	TravelTime      float64   `gorm:"type:decimal(8,2);default:0" json:"travel_time"`
	DataSource      string    `gorm:"size:50;default:sensor" json:"data_source,omitempty"`
	IsValidated     bool      `gorm:"default:true" json:"is_validated"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Location Location `gorm:"foreignKey:LocationID;references:ID" json:"location"`
}

// =====================================================
// ANALYTICS RESULTS
// =====================================================
type AnalyticsResult struct {
	ID                int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID              string    `gorm:"type:uuid;default:uuid_generate_v4()" json:"uuid"`
	AnalysisTimestamp time.Time `gorm:"autoCreateTime" json:"analysis_timestamp"`
	PeriodStart       time.Time `gorm:"not null" json:"period_start" validate:"required"`
	PeriodEnd         time.Time `gorm:"not null" json:"period_end" validate:"required"`
	LocationID        *string   `gorm:"size:50" json:"location_id,omitempty"`
	MetricType        string    `gorm:"size:50;not null" json:"metric_type" validate:"required"`
	Value             float64   `gorm:"type:decimal(15,4);not null" json:"value" validate:"required"`
	Unit              string    `gorm:"size:20" json:"unit,omitempty"`
	ConfidenceLevel   *float64  `gorm:"type:decimal(5,4)" json:"confidence_level,omitempty" validate:"omitempty,min=0,max=1"`
	Trend             string    `gorm:"size:20" json:"trend,omitempty" validate:"omitempty,oneof=increasing decreasing stable"`
	SampleSize        *int      `json:"sample_size,omitempty"`
	AggregationMethod string    `gorm:"size:50" json:"aggregation_method,omitempty"`
	IsAnomaly         bool      `gorm:"default:false" json:"is_anomaly"`
	Metadata          []byte    `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Location *Location `gorm:"foreignKey:LocationID;references:ID" json:"location,omitempty"`
}

const (
	// Congestion Levels
	CongestionLow    = "low"
	CongestionMedium = "medium"
	CongestionHigh   = "high"
	CongestionSevere = "severe"

	// Alert Types
	AlertTypeCongestion     = "congestion"
	AlertTypeAccident       = "accident"
	AlertTypeSpeedViolation = "speed_violation"
	AlertTypeSystemError    = "system_error"

	// Alert Severities
	SeverityLow      = "low"
	SeverityMedium   = "medium"
	SeverityHigh     = "high"
	SeverityCritical = "critical"

	// Alert Statuses
	AlertStatusActive       = "active"
	AlertStatusResolved     = "resolved"
	AlertStatusAcknowledged = "acknowledged"

	// Metric Types
	MetricTypeCounter   = "counter"
	MetricTypeGauge     = "gauge"
	MetricTypeHistogram = "histogram"

	// Analytics Metric Types
	MetricAvgVehicleCount = "avg_vehicle_count"
	MetricAvgSpeed        = "avg_speed"
	MetricCongestionIndex = "congestion_index"
	MetricPeakHour        = "peak_hour"
	MetricTrafficDensity  = "traffic_density"
	MetricFlowRate        = "flow_rate"
	MetricTravelTime      = "travel_time"
	MetricDelayIndex      = "delay_index"

	// Units
	UnitVehiclesPerHour = "vehicles/hour"
	UnitKmPerHour       = "km/h"
	UnitIndex           = "index"
	UnitBoolean         = "boolean"
	UnitVehiclesPerKm   = "vehicles/km"
	UnitMinutes         = "minutes"
)

// =====================================================
// ALERTS
// =====================================================
type Alert struct {
	ID               int        `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID             string     `gorm:"type:uuid;default:uuid_generate_v4()" json:"uuid"`
	Timestamp        time.Time  `gorm:"autoCreateTime" json:"timestamp"`
	LocationID       *string    `gorm:"size:50" json:"location_id,omitempty"`
	AlertType        string     `gorm:"size:50;not null" json:"alert_type" validate:"required"`
	Severity         string     `gorm:"size:20;not null" json:"severity" validate:"required,oneof=low medium high critical"`
	Message          string     `gorm:"not null" json:"message" validate:"required"`
	Description      string     `json:"description,omitempty"`
	Value            *float64   `gorm:"type:decimal(15,4)" json:"value,omitempty"`
	Threshold        *float64   `gorm:"type:decimal(15,4)" json:"threshold,omitempty"`
	Status           string     `gorm:"size:20;default:active" json:"status" validate:"oneof=active resolved acknowledged suppressed"`
	Category         string     `gorm:"size:50;default:traffic" json:"category"`
	Priority         int        `gorm:"default:1" json:"priority"`
	AssignedTo       string     `gorm:"size:100" json:"assigned_to,omitempty"`
	ResolvedAt       *time.Time `json:"resolved_at,omitempty"`
	ResolvedBy       string     `gorm:"size:100" json:"resolved_by,omitempty"`
	ResolutionNotes  string     `json:"resolution_notes,omitempty"`
	NotificationSent bool       `gorm:"default:false" json:"notification_sent"`
	Metadata         []byte     `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	Location *Location `gorm:"foreignKey:LocationID;references:ID" json:"location,omitempty"`
}

// =====================================================
// SYSTEM METRICS
// =====================================================
type SystemMetric struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID        string    `gorm:"type:uuid;default:uuid_generate_v4()" json:"uuid"`
	ServiceName string    `gorm:"size:100;not null" json:"service_name" validate:"required"`
	MetricName  string    `gorm:"size:100;not null" json:"metric_name" validate:"required"`
	MetricType  string    `gorm:"size:20;not null" json:"metric_type" validate:"required,oneof=counter gauge histogram summary"`
	Value       float64   `gorm:"type:decimal(15,4);not null" json:"value" validate:"required"`
	Labels      []byte    `gorm:"type:jsonb" json:"labels,omitempty"`
	Timestamp   time.Time `gorm:"autoCreateTime" json:"timestamp"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// =====================================================
// CONFIGURATIONS
// =====================================================
type Configuration struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID        string    `gorm:"type:uuid;default:uuid_generate_v4()" json:"uuid"`
	Key         string    `gorm:"size:255;not null;unique" json:"key" validate:"required"`
	Value       string    `gorm:"not null" json:"value" validate:"required"`
	Description string    `json:"description,omitempty"`
	Category    string    `gorm:"size:100" json:"category,omitempty"`
	DataType    string    `gorm:"size:50;default:string" json:"data_type,omitempty"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	IsEncrypted bool      `gorm:"default:false" json:"is_encrypted"`
	CreatedBy   string    `gorm:"size:100" json:"created_by,omitempty"`
	UpdatedBy   string    `gorm:"size:100" json:"updated_by,omitempty"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// =====================================================
// AUDIT LOGS
// =====================================================
type AuditLog struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID      string    `gorm:"type:uuid;default:uuid_generate_v4()" json:"uuid"`
	Action    string    `gorm:"size:100;not null" json:"action" validate:"required"`
	TableName string    `gorm:"size:100" json:"table_name,omitempty"`
	RecordID  string    `gorm:"size:100" json:"record_id,omitempty"`
	OldValues []byte    `gorm:"type:jsonb" json:"old_values,omitempty"`
	NewValues []byte    `gorm:"type:jsonb" json:"new_values,omitempty"`
	UserID    string    `gorm:"size:100" json:"user_id,omitempty"`
	IPAddress string    `gorm:"size:45" json:"ip_address,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
	Timestamp time.Time `gorm:"autoCreateTime" json:"timestamp"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
