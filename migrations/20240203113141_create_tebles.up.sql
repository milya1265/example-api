CREATE TABLE IF NOT EXISTS users
(
    id  varchar(40) unique PRIMARY KEY,
    login varchar(  30) unique,
    password_hash varchar(100)
);

CREATE TABLE IF NOT EXISTS jwt
(
    user_id           varchar(40) unique  not null ,
    foreign key (user_id) REFERENCES users (id),
    access_token varchar(200)  not null,
    refresh_token varchar(200) not null
);

CREATE TABLE IF NOT EXISTS users_role
(
    user_id   varchar(40) unique,
    foreign key (user_id) REFERENCES users (id),
    role integer
);

CREATE TABLE IF NOT EXISTS  warehouse
(
    id        serial primary key,
    name      VARCHAR(100) unique,
    available BOOLEAN
);

CREATE TABLE IF NOT EXISTS product
(
    id          serial primary key,
    name        VARCHAR(100),
    size        INTEGER,
    unique_code VARCHAR(100) unique
);

CREATE TABLE IF NOT EXISTS warehouse_product
(
    id           serial primary key,
    warehouse_id INTEGER REFERENCES warehouse (id),
    product_code VARCHAR(100) REFERENCES product (unique_code),
    total_count  INTEGER,
    left_count   INTEGER
);

CREATE TABLE IF NOT EXISTS reservation
(
    id           serial primary key,
    warehouse_id INTEGER REFERENCES warehouse (id),
    product_code VARCHAR(100) REFERENCES product (unique_code),
    count        INTEGER
);