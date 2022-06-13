CREATE TABLE users (
    id bigserial not null,
    email varchar not null unique,
    encrypted_password varchar not null
);
