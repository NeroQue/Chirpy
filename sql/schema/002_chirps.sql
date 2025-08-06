-- +goose Up
CREATE TABLE chirps (
    id UUID primary key,
    created_at TIMESTAMP not null,
    updated_at TIMESTAMP not null,
    body text not null,
    user_id UUID not null references users(id) on delete cascade
);

-- +goose Down
DROP TABLE chirps;