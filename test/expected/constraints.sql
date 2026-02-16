CREATE TABLE persons (
    id BIGINT PRIMARY KEY NOT NULL,
    name VARCHAR(255) CHECK (length(name) > 0),
    age INTEGER CHECK (age >= 18),
    email VARCHAR(255) CHECK (email ~* '^[a-z0-9._%+-]+@[a-z0-9.-]+\\.[a-z]{2,}$'),
    active BOOLEAN NOT NULL DEFAULT true,
    role VARCHAR(255) DEFAULT 'user',
    created BIGINT NOT NULL DEFAULT extract(epoch from now())
);
