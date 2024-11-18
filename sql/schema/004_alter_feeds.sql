-- +goose Up
ALTER TABLE feeds
ADD last_fechted_at TIMESTAMP;

-- +goose Down
ALTER TABLE feeds
DROP COLUMN last_fechted_at;