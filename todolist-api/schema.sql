-- file: schema.sql
CREATE TABLE tasks (
    id    SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    completed  BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
