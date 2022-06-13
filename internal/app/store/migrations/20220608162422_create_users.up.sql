CREATE TABLE users (
    id bigserial not null,
    email varchar not null unique,
    encrypted_password varchar not null
);

CREATE TABLE transactions (
    trans_id bigserial not null,
    user_id bigint not null,
    pay numeric not null,
    currency varchar not null,
    time_create timestamp not null,
    time_update timestamp not null,
    trans_status varchar not null
);