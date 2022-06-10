CREATE TABLE users (
    transaction_id bigserial not null primary key,
    id bigserial not null,
    email varchar not null,
    pay numeric,
    currency varchar,
    time_create time not null,
    time_update time not null,
    transaction_status varchar,
    encrypted_password varchar
);
