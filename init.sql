DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'glifdb') THEN
        CREATE DATABASE glifdb;
ELSE
        RAISE NOTICE 'Database "glifdb" already exists, skipping creation.';
END IF;
END $$;

DO $$
BEGIN
    CREATE USER glif_user WITH ENCRYPTED PASSWORD 'qwerty12345678';
EXCEPTION
    WHEN duplicate_object THEN
        RAISE NOTICE 'user glif_user already exists, skipping';
END $$;

GRANT ALL PRIVILEGES ON DATABASE glifdb TO glif_user;

\c glifdb;

CREATE TABLE transactions (
                              tx_hash TEXT PRIMARY KEY,
                              sender TEXT NOT NULL,
                              receiver TEXT NOT NULL,
                              amount NUMERIC NOT NULL,
                              status TEXT NOT NULL,
                              timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_sender ON transactions(sender);
CREATE INDEX idx_receiver ON transactions(receiver);