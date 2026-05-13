CREATE TABLE IF NOT EXISTS inventories (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL,
    stock INT NOT NULL
);