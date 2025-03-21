-- Create tables for storing system information
CREATE TABLE IF NOT EXISTS system_info (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    os_name VARCHAR(255) NOT NULL,
    os_version VARCHAR(255) NOT NULL,
    os_platform VARCHAR(255) NOT NULL,
    osquery_version VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS installed_apps (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    system_info_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    path VARCHAR(512) NOT NULL,
    bundle_identifier VARCHAR(255),
    bundle_name VARCHAR(255),
    bundle_short_version VARCHAR(100),
    display_name VARCHAR(255),
    minimum_system_version VARCHAR(50),
    last_opened_time DOUBLE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    end_time TIMESTAMP NULL DEFAULT NULL,
    FOREIGN KEY (system_info_id) REFERENCES system_info(id) ON DELETE CASCADE
);

-- Indexes for better query performance
CREATE INDEX idx_system_info_created_at ON system_info(created_at);
CREATE INDEX idx_installed_apps_system_info_id ON installed_apps(system_info_id);
CREATE INDEX idx_installed_apps_bundle_identifier ON installed_apps(bundle_identifier);
CREATE INDEX idx_installed_apps_end_time ON installed_apps(end_time); 