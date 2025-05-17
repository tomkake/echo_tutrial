
-- +migrate Up
CREATE TABLE Users(
    id binary(16) PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255),
    password VARCHAR(60),
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) COMMENT "ユーザーテーブル";

-- +migrate Down
DROP TABLE users;