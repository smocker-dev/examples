-- +migrate Down

DROP SCHEMA public CASCADE;
CREATE SCHEMA public;

-- +migrate Up
CREATE TABLE reservations (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    hotel_id BIGINT NOT NULL,
    room_number BIGINT NOT NULL
);