-- +goose Up
CREATE TABLE refresh_tokens (
    token TEXT PRIMARY KEY,
    created_at TIMESTAMP not null,
    updated_at TIMESTAMP not null,
    user_id UUID not null references users(id) on delete cascade,
    expires_at TIMESTAMP not null,
    revoked_at TIMESTAMP default NULL
);

-- +goose Down
drop table refresh_tokens;