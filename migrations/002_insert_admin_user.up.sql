INSERT INTO users (email, password, name, role, created_at, updated_at) 
VALUES (
    'admin@rumi.play',
    '$2a$10$f6tPbpGU9eOC6vMTK3Nz4e.4Qs9OvjIgXrATG4Asc/6QZlhcsI0VK', -- adminrumi
    'Admin User',
    'admin',
    NOW(),
    NOW()
) ON DUPLICATE KEY UPDATE
    password = VALUES(password),
    name = VALUES(name),
    role = VALUES(role),
    updated_at = NOW();
