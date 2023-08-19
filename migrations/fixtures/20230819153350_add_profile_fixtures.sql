-- +goose Up
-- +goose StatementBegin
INSERT INTO profiles (id, user_id, first_name, last_name, email, phone, dob)
VALUES
    ('a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'John', 'Doe', 'johndoe@example.com', '+380950001122', NULL),
    ('a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'Jane', 'Doe', 'janedoe@example.com', '+380950001123', '1990-12-21'),
    ('a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'Jack', 'Doe', 'jackdoe@example.com', '+380950001124', NULL),
    ('a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 'Jill', '', 'jilldoe@example.com', '+380950001125', '1991-01-01'),
    ('a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', 'John', 'Smith', 'smith@test.com', '+380950001126', CURRENT_DATE)
ON CONFLICT (user_id) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM profiles
WHERE id IN ('a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 'a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a15');
-- +goose StatementEnd
