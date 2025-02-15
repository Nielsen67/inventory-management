CREATE DATABASE inventory_management;

DROP DATABASE inventory_management;

USE inventory_management;

INSERT INTO products (name, description, price, category, created_at, updated_at)
VALUES
    ('Ring of Dagan', 'Ring that possess the power of Dagnue', 1200000000, 'Accessories', NOW(), NOW()),
    ('Ring of Dagnue', 'Ring that possess the power of Dagan', 1100000000, 'Accessories', NOW(), NOW()),
    ('Dagon Leash', 'Necklace that possess the power of Dagon', 1100000000, 'Accessories', NOW(), NOW());



INSERT INTO inventories (product_id, stock, location, created_at, updated_at)
VALUES
    (9, 1, 'Cengkareng', NOW(), NOW()),
    (10, 2, 'Pluit', NOW(), NOW()),
    (11, 2, 'Cakung', NOW(), NOW());


INSERT INTO orders (date, customer_name, shipping_address, total)
VALUES
    (NOW(), 'Arthur', 'Jakarta', 1200000000);
    
INSERT INTO order_details (product_id, order_id, qty, price_detail)
VALUES
    (9, 10, 1, 1200000000);


-- PRODUCTS & INVENTORIES
SELECT p.name, p.description, i.stock, i.location
FROM products p
JOIN inventories i ON p.id = i.product_id;

-- ORDERS & ORDER DETAILS
SELECT o.date, o.customer_name, od.qty, od.price_detail
FROM orders o
JOIN order_details od ON o.id = od.order_id;

-- Total Pesanan per Produk
SELECT p.name, SUM(od.qty * od.price_detail) AS total_pesanan
FROM products p
JOIN order_details od ON p.id = od.product_id
GROUP BY p.name;

-- Tingkat Stok per Lokasi
SELECT location, product_id, SUM(stock) AS total_stok
FROM inventories
GROUP BY location, product_id;
