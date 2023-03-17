CREATE TABLE "url" (
    token TEXT PRIMARY KEY,
    target_url TEXT NOT NULL,
    visits INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP
);