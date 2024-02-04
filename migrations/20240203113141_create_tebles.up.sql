CREATE TABLE warehouse
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