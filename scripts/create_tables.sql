-- =====================================================
-- TRAFFIC ANALYTICS DATABASE SCHEMA
-- PostgreSQL Database Creation Script
-- =====================================================

-- =====================================================
-- EXTENSIONS
-- =====================================================
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =====================================================
-- LOCATIONS TABLE
-- Stores information about traffic monitoring locations
-- =====================================================
CREATE TABLE locations (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    address TEXT,
    city VARCHAR(100),
    country VARCHAR(100) DEFAULT 'Argentina',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for locations
CREATE INDEX idx_locations_coordinates ON locations(latitude, longitude);
CREATE INDEX idx_locations_city ON locations(city);
CREATE INDEX idx_locations_active ON locations(is_active);

-- =====================================================
-- TRAFFIC DATA TABLE
-- Stores raw real-time traffic data
-- =====================================================
CREATE TABLE traffic_data (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4(),
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    location_id VARCHAR(50) NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    vehicle_count INTEGER NOT NULL CHECK (vehicle_count >= 0),
    average_speed DECIMAL(5, 2) NOT NULL CHECK (average_speed >= 0),
    congestion_level VARCHAR(20) NOT NULL CHECK (congestion_level IN ('low', 'medium', 'high', 'severe')),
    max_speed DECIMAL(5, 2) CHECK (max_speed >= 0),
    min_speed DECIMAL(5, 2) CHECK (min_speed >= 0),
    occupancy DECIMAL(5, 2) CHECK (occupancy >= 0 AND occupancy <= 100),
    queue_length DECIMAL(8, 2) DEFAULT 0,
    travel_time DECIMAL(8, 2) DEFAULT 0,
    data_source VARCHAR(50) DEFAULT 'sensor',
    is_validated BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for traffic_data
CREATE INDEX idx_traffic_data_timestamp ON traffic_data(timestamp);
CREATE INDEX idx_traffic_data_location ON traffic_data(location_id);
CREATE INDEX idx_traffic_data_congestion ON traffic_data(congestion_level);
CREATE INDEX idx_traffic_data_location_time ON traffic_data(location_id, timestamp);
CREATE INDEX idx_traffic_data_uuid ON traffic_data(uuid);

-- =====================================================
-- ANALYTICS RESULTS TABLE
-- Stores processed analytics results
-- =====================================================
CREATE TABLE analytics_results (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4(),
    analysis_timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    period_start TIMESTAMP NOT NULL,
    period_end TIMESTAMP NOT NULL,
    location_id VARCHAR(50) REFERENCES locations(id) ON DELETE CASCADE,
    metric_type VARCHAR(50) NOT NULL,
    value DECIMAL(15, 4) NOT NULL,
    unit VARCHAR(20),
    confidence_level DECIMAL(5, 4) CHECK (confidence_level >= 0 AND confidence_level <= 1),
    trend VARCHAR(20) CHECK (trend IN ('increasing', 'decreasing', 'stable')),
    sample_size INTEGER,
    aggregation_method VARCHAR(50),
    is_anomaly BOOLEAN DEFAULT FALSE,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for analytics_results
CREATE INDEX idx_analytics_results_timestamp ON analytics_results(analysis_timestamp);
CREATE INDEX idx_analytics_results_location ON analytics_results(location_id);
CREATE INDEX idx_analytics_results_metric_type ON analytics_results(metric_type);
CREATE INDEX idx_analytics_results_period ON analytics_results(period_start, period_end);
CREATE INDEX idx_analytics_results_anomaly ON analytics_results(is_anomaly);
CREATE INDEX idx_analytics_results_uuid ON analytics_results(uuid);

-- =====================================================
-- ALERTS TABLE
-- Stores system and traffic alerts
-- =====================================================
CREATE TABLE alerts (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4(),
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    location_id VARCHAR(50) REFERENCES locations(id) ON DELETE CASCADE,
    alert_type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    message TEXT NOT NULL,
    description TEXT,
    value DECIMAL(15, 4),
    threshold DECIMAL(15, 4),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'resolved', 'acknowledged', 'suppressed')),
    category VARCHAR(50) DEFAULT 'traffic',
    priority INTEGER DEFAULT 1,
    assigned_to VARCHAR(100),
    resolved_at TIMESTAMP,
    resolved_by VARCHAR(100),
    resolution_notes TEXT,
    notification_sent BOOLEAN DEFAULT FALSE,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for alerts
CREATE INDEX idx_alerts_timestamp ON alerts(timestamp);
CREATE INDEX idx_alerts_location ON alerts(location_id);
CREATE INDEX idx_alerts_type ON alerts(alert_type);
CREATE INDEX idx_alerts_severity ON alerts(severity);
CREATE INDEX idx_alerts_status ON alerts(status);
CREATE INDEX idx_alerts_category ON alerts(category);
CREATE INDEX idx_alerts_active ON alerts(status) WHERE status = 'active';
CREATE INDEX idx_alerts_uuid ON alerts(uuid);

-- =====================================================
-- SYSTEM METRICS TABLE
-- Stores system performance metrics for monitoring
-- =====================================================
CREATE TABLE system_metrics (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4(),
    service_name VARCHAR(100) NOT NULL,
    metric_name VARCHAR(100) NOT NULL,
    metric_type VARCHAR(20) NOT NULL CHECK (metric_type IN ('counter', 'gauge', 'histogram', 'summary')),
    value DECIMAL(15, 4) NOT NULL,
    labels JSONB,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for system_metrics
CREATE INDEX idx_system_metrics_service ON system_metrics(service_name);
CREATE INDEX idx_system_metrics_name ON system_metrics(metric_name);
CREATE INDEX idx_system_metrics_timestamp ON system_metrics(timestamp);
CREATE INDEX idx_system_metrics_labels ON system_metrics USING GIN (labels);
CREATE INDEX idx_system_metrics_uuid ON system_metrics(uuid);

-- =====================================================
-- CONFIGURATION TABLE
-- Stores system configuration parameters
-- =====================================================
CREATE TABLE configurations (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4(),
    key VARCHAR(255) NOT NULL UNIQUE,
    value TEXT NOT NULL,
    description TEXT,
    category VARCHAR(100),
    data_type VARCHAR(50) DEFAULT 'string',
    is_active BOOLEAN DEFAULT TRUE,
    is_encrypted BOOLEAN DEFAULT FALSE,
    created_by VARCHAR(100),
    updated_by VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for configurations
CREATE INDEX idx_configurations_key ON configurations(key);
CREATE INDEX idx_configurations_category ON configurations(category);
CREATE INDEX idx_configurations_active ON configurations(is_active);
CREATE INDEX idx_configurations_uuid ON configurations(uuid);

-- =====================================================
-- AUDIT LOG TABLE
-- Stores audit trail of system operations
-- =====================================================
CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4(),
    action VARCHAR(100) NOT NULL,
    table_name VARCHAR(100),
    record_id VARCHAR(100),
    old_values JSONB,
    new_values JSONB,
    user_id VARCHAR(100),
    ip_address VARCHAR(45),
    user_agent TEXT,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for audit_logs
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_table ON audit_logs(table_name);
CREATE INDEX idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_timestamp ON audit_logs(timestamp);
CREATE INDEX idx_audit_logs_uuid ON audit_logs(uuid);

-- =====================================================
-- TRIGGERS FOR UPDATED_AT FIELDS
-- =====================================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for tables with updated_at
CREATE TRIGGER update_locations_updated_at 
    BEFORE UPDATE ON locations 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_traffic_data_updated_at 
    BEFORE UPDATE ON traffic_data 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_analytics_results_updated_at 
    BEFORE UPDATE ON analytics_results 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_alerts_updated_at 
    BEFORE UPDATE ON alerts 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_configurations_updated_at 
    BEFORE UPDATE ON configurations 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- INITIAL DATA INSERTION
-- =====================================================

-- Insert default configurations
INSERT INTO configurations (key, value, description, category, data_type) VALUES
('system.version', '1.0.0', 'Current system version', 'system', 'string'),
('analytics.processing_interval', '300', 'Analytics processing interval in seconds', 'analytics', 'integer'),
('alert.congestion_threshold', '0.7', 'Congestion threshold for alerts', 'alerts', 'decimal'),
('alert.speed_threshold', '15.0', 'Minimum speed threshold for accident detection', 'alerts', 'decimal'),
('cache.ttl_default', '3600', 'Default cache TTL in seconds', 'cache', 'integer'),
('kafka.broker', 'kafka:9092', 'Kafka broker address', 'kafka', 'string'),
('redis.host', 'redis:6379', 'Redis host address', 'redis', 'string'),
('database.max_connections', '100', 'Maximum database connections', 'database', 'integer');

-- =====================================================
-- VIEWS FOR COMMON QUERIES
-- =====================================================

-- View for current traffic status
CREATE VIEW current_traffic_status AS
SELECT 
    td.location_id,
    l.name as location_name,
    td.timestamp,
    td.vehicle_count,
    td.average_speed,
    td.congestion_level,
    td.occupancy,
    td.travel_time
FROM traffic_data td
JOIN locations l ON td.location_id = l.id
WHERE td.timestamp >= NOW() - INTERVAL '30 minutes'
ORDER BY td.timestamp DESC;

-- View for active alerts
CREATE VIEW active_alerts AS
SELECT 
    a.id,
    a.timestamp,
    a.location_id,
    l.name as location_name,
    a.alert_type,
    a.severity,
    a.message,
    a.value,
    a.threshold
FROM alerts a
LEFT JOIN locations l ON a.location_id = l.id
WHERE a.status = 'active'
ORDER BY a.severity DESC, a.timestamp DESC;

-- View for analytics summary
CREATE VIEW analytics_summary AS
SELECT 
    ar.location_id,
    l.name as location_name,
    ar.metric_type,
    AVG(ar.value) as avg_value,
    MAX(ar.value) as max_value,
    MIN(ar.value) as min_value,
    COUNT(*) as sample_count,
    MAX(ar.analysis_timestamp) as last_analysis
FROM analytics_results ar
JOIN locations l ON ar.location_id = l.id
WHERE ar.analysis_timestamp >= NOW() - INTERVAL '24 hours'
GROUP BY ar.location_id, l.name, ar.metric_type;

-- =====================================================
-- INDEXES SUMMARY
-- =====================================================
/*
Created indexes:
- locations: coordinates, city, active status
- traffic_data: timestamp, location, congestion, composite indexes
- analytics_results: timestamp, location, metric type, period, anomalies
- alerts: timestamp, location, type, severity, status, active alerts
- system_metrics: service, metric name, timestamp, labels
- configurations: key, category, active status
- audit_logs: action, table, user, timestamp
*/

-- =====================================================
-- DATABASE SCHEMA COMPLETE
-- =====================================================