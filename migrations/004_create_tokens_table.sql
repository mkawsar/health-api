-- +goose Up
-- Create tokens table
CREATE TABLE IF NOT EXISTS tokens (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    user_id INT NOT NULL,
    token TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    expries_at DATETIME NOT NULL,
    blacklisted TINYINT(1) NOT NULL DEFAULT 0,
    INDEX idx_tokens_user_id (user_id),
    INDEX idx_tokens_token (token(255)),
    INDEX idx_tokens_expries_at (expries_at),
    INDEX idx_tokens_deleted_at (deleted_at),
    CONSTRAINT fk_tokens_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS tokens;

