CREATE TABLE IF NOT EXISTS contracts (
                                         id BIGINT PRIMARY KEY,
                                         data JSONB
);

CREATE TABLE IF NOT EXISTS states (
                                      id BIGINT PRIMARY KEY,
                                      data JSONB
);