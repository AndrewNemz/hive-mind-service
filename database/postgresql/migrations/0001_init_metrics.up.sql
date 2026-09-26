CREATE TABLE IF NOT EXISTS metrics (
    id TEXT NOT NULL,
    type TEXT NOT NULL,
    value DOUBLE PRECISION,
    delta BIGINT,
    PRIMARY KEY (id, type)
);

CREATE INDEX IF NOT EXISTS idx_metrics_type ON metrics(type);