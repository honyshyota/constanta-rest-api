CREATE TYPE pay AS (
    transaction_id integer,
    summ numeric,
    currency varchar,
    transaction_status varchar
);

CREATE TABLE users (
    user_id bigserial not null,
    email varchar not null,
    pay pay[],
    time_create time not null,
    time_update time not null
);