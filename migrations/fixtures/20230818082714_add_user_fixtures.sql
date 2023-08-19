-- +goose Up
-- +goose StatementBegin
INSERT INTO users (id, login, role, password_hash)
VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'user1', 'ROLE_ADMIN', '$2a$10$z4FTAQvrdm5mqV0hDLxH/.hxEG6a1gaydnhqPF/COB04QdX2CDtPe'),
    ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'user2', 'ROLE_ADMIN', '$2a$10$z2PjrBRBQoQqKUu44OPJ0udLFjTZhXFtmZsuBGUE4QXKVRzYZh/T6'),
    ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'user3', 'ROLE_USER', '$2a$10$wGbyOA9vB6tXgtGrojae3u2SIZLIn7K40UDQC7AKGt7A/WWIS2TDG'),
    ('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 'user4', 'ROLE_USER', '$2a$10$T6dZdTjkobi9YudBuNZmfOaZpeJPqL4bUBmW/0SjnaMzlDtn1wEkq'),
    ('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', 'user5', 'ROLE_USER', '$2a$10$a4kxrWgN2MCoWKTCGyILSOStnGwsjGaL0lY7FnDV10r2enqTPWV1O')
ON CONFLICT (login)
    DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users
WHERE login IN ('user1', 'user2', 'user3', 'user4', 'user5');
-- +goose StatementEnd
