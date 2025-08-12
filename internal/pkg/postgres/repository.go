package postgres

import (
	"api-traffic-analytics/internal/shared/models"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// TrafficDataRepository handles CRUD operations for TrafficData using GORM
type TrafficDataRepository struct {
	db *gorm.DB
}

// NewTrafficDataRepository creates a new instance of TrafficDataRepository
func NewTrafficDataRepository(db *gorm.DB) *TrafficDataRepository {
	return &TrafficDataRepository{db: db}
}

// Create inserts a new traffic data record
func (r *TrafficDataRepository) Create(ctx context.Context, data *models.TrafficData) (*models.TrafficData, error) {

	// Use GORM to create the record
	result := r.db.WithContext(ctx).Create(data)
	if result.Error != nil {
		return nil, fmt.Errorf("error creating traffic data: %w", result.Error)
	}

	return data, nil
}

// GetByID retrieves a traffic data record by its ID
func (r *TrafficDataRepository) GetByID(ctx context.Context, id int64) (*models.TrafficData, error) {
	var data models.TrafficData

	result := r.db.WithContext(ctx).Where("id = ?", id).First(&data)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("traffic data not found with id %d", id)
		}
		return nil, fmt.Errorf("error getting traffic data by ID: %w", result.Error)
	}

	return &data, nil
}

// Update modifies an existing traffic data record
func (r *TrafficDataRepository) Update(ctx context.Context, data *models.TrafficData) error {
	// Use GORM to update the record
	result := r.db.WithContext(ctx).Save(data)
	if result.Error != nil {
		return fmt.Errorf("error updating traffic data: %w", result.Error)
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return fmt.Errorf("no traffic data found with id %d to update", data)
	}

	return nil
}

// Delete removes a traffic data record by its ID
func (r *TrafficDataRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.TrafficData{})
	if result.Error != nil {
		return fmt.Errorf("error deleting traffic data: %w", result.Error)
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return fmt.Errorf("no traffic data found with id %d to delete", id)
	}

	return nil
}

// GetAll retrieves all traffic data records with optional pagination
func (r *TrafficDataRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.TrafficData, error) {
	var data []*models.TrafficData

	query := r.db.WithContext(ctx).Order("timestamp DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	result := query.Find(&data)
	if result.Error != nil {
		return nil, fmt.Errorf("error getting all traffic data: %w", result.Error)
	}

	return data, nil
}

// GetByLocation retrieves traffic data records for a specific location
func (r *TrafficDataRepository) GetByLocation(ctx context.Context, locationID string, limit int) ([]*models.TrafficData, error) {
	var data []*models.TrafficData

	result := r.db.WithContext(ctx).
		Where("location_id = ?", locationID).
		Order("timestamp DESC").
		Limit(limit).
		Find(&data)

	if result.Error != nil {
		return nil, fmt.Errorf("error getting traffic data by location: %w", result.Error)
	}

	return data, nil
}

// GetByTimeRange retrieves traffic data within a specific time range
func (r *TrafficDataRepository) GetByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*models.TrafficData, error) {
	var data []*models.TrafficData

	result := r.db.WithContext(ctx).
		Where("timestamp >= ? AND timestamp <= ?", startTime, endTime).
		Order("timestamp ASC").
		Find(&data)

	if result.Error != nil {
		return nil, fmt.Errorf("error getting traffic data by time range: %w", result.Error)
	}

	return data, nil
}

// GetByLocationAndTimeRange retrieves traffic data for a location within a time range
func (r *TrafficDataRepository) GetByLocationAndTimeRange(ctx context.Context, locationID string, startTime, endTime time.Time) ([]*models.TrafficData, error) {
	var data []*models.TrafficData

	result := r.db.WithContext(ctx).
		Where("location_id = ? AND timestamp >= ? AND timestamp <= ?", locationID, startTime, endTime).
		Order("timestamp ASC").
		Find(&data)

	if result.Error != nil {
		return nil, fmt.Errorf("error getting traffic data by location and time range: %w", result.Error)
	}

	return data, nil
}

// GetLatestByLocation retrieves the most recent traffic data for a location
func (r *TrafficDataRepository) GetLatestByLocation(ctx context.Context, locationID string) (*models.TrafficData, error) {
	var data models.TrafficData

	result := r.db.WithContext(ctx).
		Where("location_id = ?", locationID).
		Order("timestamp DESC").
		First(&data)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no traffic data found for location %s", locationID)
		}
		return nil, fmt.Errorf("error getting latest traffic data by location: %w", result.Error)
	}

	return &data, nil
}

// BatchCreate inserts multiple traffic data records efficiently
func (r *TrafficDataRepository) BatchCreate(ctx context.Context, dataList []*models.TrafficData) error {
	// Use GORM's CreateInBatches for efficient batch insertion
	result := r.db.WithContext(ctx).CreateInBatches(dataList, 100)
	if result.Error != nil {
		return fmt.Errorf("error creating traffic data batch: %w", result.Error)
	}

	return nil
}

// Count returns the total number of traffic data records
func (r *TrafficDataRepository) Count(ctx context.Context) (int64, error) {
	var count int64

	result := r.db.WithContext(ctx).Model(&models.TrafficData{}).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("error counting traffic data records: %w", result.Error)
	}

	return count, nil
}

// CountByLocation returns the number of traffic data records for a specific location
func (r *TrafficDataRepository) CountByLocation(ctx context.Context, locationID string) (int64, error) {
	var count int64

	result := r.db.WithContext(ctx).
		Model(&models.TrafficData{}).
		Where("location_id = ?", locationID).
		Count(&count)

	if result.Error != nil {
		return 0, fmt.Errorf("error counting traffic data records by location: %w", result.Error)
	}

	return count, nil
}

func (r *TrafficDataRepository) CreateAnalyticsResult(ctx context.Context, result *models.AnalyticsResult) error {
	return r.db.WithContext(ctx).Create(result).Error
}
