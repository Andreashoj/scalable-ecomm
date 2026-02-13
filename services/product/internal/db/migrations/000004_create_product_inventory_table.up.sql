CREATE TABLE product_inventory
(
    product_id UUID,
    amount     INT DEFAULT 0,
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (product_id) REFERENCES product (id) ON DELETE CASCADE
);