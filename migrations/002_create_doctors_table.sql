-- +goose Up
-- Create doctors table
CREATE TABLE IF NOT EXISTS doctors (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    name VARCHAR(255) NOT NULL,
    specialization VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    experience TEXT,
    location VARCHAR(255) NOT NULL,
    license VARCHAR(255) NOT NULL UNIQUE,
    work_hours VARCHAR(100),
    availability TINYINT(1) NOT NULL DEFAULT 1,
    work_days JSON,
    work_time JSON,
    work_time_end JSON,
    UNIQUE INDEX idx_doctors_license (license),
    INDEX idx_doctors_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS doctors;

