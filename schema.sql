CREATE TYPE file_type_enum AS ENUM ('LOCAL', 'STREAMING');

CREATE TABLE service_logs (
    id BIGSERIAL PRIMARY KEY,
    service_name VARCHAR(255) NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE file_downloads (
    id BIGSERIAL PRIMARY KEY,
    file_path VARCHAR NOT NULL,
    file_type file_type_enum NOT NULL,
    expiry INTEGER NOT NULL,
    client_key VARCHAR(255) NOT NULL,
    server_key VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);
